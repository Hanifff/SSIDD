using System;
using System.Threading;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Http;
using System.Collections.Generic;
using System.IO;
using System.Text;
using Newtonsoft.Json;
using Newtonsoft.Json.Linq;
using System.Linq;
using Microsoft.AspNetCore.Mvc;
using Hyperledger.Aries.Agents;
using Hyperledger.Aries.Configuration;
using Hyperledger.Aries.Extensions;
using Hyperledger.Aries.Features.PresentProof;
using Hyperledger.Aries.Features.Handshakes.Connection;
using Hyperledger.Aries.Features.Handshakes.Connection.Models;
using Hyperledger.Aries.Models.Events;
using Hyperledger.Aries.Features.Handshakes.DidExchange;
using Hyperledger.Aries.Features.Handshakes.Common;
using Hyperledger.Aries.Features.IssueCredential;
using Hyperledger.Aries.Features.PresentProof.Messages;
using Hyperledger.Aries.Models.Records;
using Hyperledger.Aries.Storage;
using Hyperledger.Indy.DidApi;
using Hyperledger.Indy.LedgerApi;
using WebAgent.Models;
using Hyperledger.Indy.AnonCredsApi;
using System.Diagnostics;
using WebAgent.Utils;
using Grpc.Net.Client;
using Google.Protobuf.WellKnownTypes;
using Google.Protobuf;
namespace WebAgent.Controllers
{
    public class AccessController : Controller
    {
        private readonly IAgentProvider _agentContextProvider;
        private readonly IProvisioningService _provisionService;
        private readonly IWalletService _walletService;
        private readonly IConnectionService _connectionService;
        private readonly IProofService _proofService;
        private readonly IMessageService _messageService;

        private Requests _requests;
        public AccessController(
            IAgentProvider agentContextProvider,
            IProvisioningService provisionService,
            IWalletService walletService,
            IConnectionService connectionService,
            IProofService proofService,
            IMessageService messageService,
            Requests requests)
        {
            _agentContextProvider = agentContextProvider;
            _provisionService = provisionService;
            _walletService = walletService;
            _connectionService = connectionService;
            _proofService = proofService;
            _messageService = messageService;
            _requests = requests;
        }

        [HttpGet]
        public async Task<IActionResult> Index()
        {
            var models = new List<RequestorModel>();
            foreach (var accreq in _requests.GetRequests())
            {
                models.Add(accreq.Value);
            }
            return View(new RequestorsModel { accessRequests = models });
        }


        [HttpPost]
        public async Task<IActionResult> AccessRequest([FromBody] AccessRequestModel model)
        {
            try
            {
                Console.WriteLine("******Recievied post request for ResourceId:  {0}", model.ResourceId);
                PolicyModel policy = FetchPolicy(model.Type);
                // TODO: should request for presentation based on policy retrieved from PAPSC smart contract 
                var context = await _agentContextProvider.GetContextAsync();
                var connections = await _connectionService.ListAsync(context);
                foreach (var connection in connections)
                {
                    if (connection.TheirDid == model.ClientDID)
                    {
                        // Verifier sends a proof request to prover
                        var proofRequestObject = new ProofRequest
                        {
                            Name = "ProofReq",
                            Version = "1.0",
                            Nonce = await AnonCreds.GenerateNonceAsync(),
                            RequestedAttributes = new Dictionary<string, ProofAttributeInfo>
                    {
                        // TODO: Requested attrs. must be based on the policy
                        {"access-requirement", new ProofAttributeInfo {Names = new [] {"status" ,"organisationid", "expiration", "nftid" } }}
                    }
                        };
                        var (proofReq, _) = await _proofService.CreateRequestAsync(context, proofRequestObject, connection.Id);
                        await _messageService.SendAsync(context, proofReq, connection);
                        // store in the request for futher handling
                        // key is the client DID
                        _requests.Add(connection.TheirDid, new RequestorModel()
                        {
                            ClientDid = connection.TheirDid,
                            ClientVK = connection.TheirVk,
                            Timestamp = connection.CreatedAtUtc ?? DateTime.MinValue,
                            RequestId = model.ConnectionId, // Matches the request id of "ProofPresentationModel object" when recieveing proof
                            Status = "VP",
                            ClientRequest = new AccessRequestModel()
                            {
                                ResourceId = model.ResourceId,
                                ConnectionId = model.ConnectionId,
                                Type = model.Type,
                                Data = model.Data,
                                Policy = model.Policy,
                                ResAttributes = model.ResAttributes,
                                Status = "VP"
                            }
                        });
                        return Ok(_requests.Get(connection.TheirDid).ClientRequest);
                    }
                }
                return StatusCode(500, new AccessRequestModel()
                {
                    ResourceId = model.ResourceId,
                    ConnectionId = model.ConnectionId,
                    Type = model.Type,
                    Data = model.Data,
                    Policy = model.Policy,
                    ResAttributes = model.ResAttributes,
                    Status = "Failed"
                });
            }
            catch (Exception e)
            {
                Console.WriteLine("\nException Caught!");
                Console.WriteLine("Message: {0} ", e.Message);
                var st = new StackTrace(e, true);
                var line = st.GetFrame(st.FrameCount - 1).GetFileLineNumber();
                Console.WriteLine("Line number: {0} ", line);
                return null;
            }
        }


        [HttpPost]
        public async Task<IActionResult> ClientProofPresentation([FromBody] ProofPresentationModel model)
        {
            try
            {
                // finding the proof record or connection. 
                var context = await _agentContextProvider.GetContextAsync();
                var accesptedProofs = await _proofService.ListAcceptedAsync(context);
                foreach (var accepted in accesptedProofs)
                {
                    // extract revealed client attribute from the proof presentation
                    JToken jToken = JToken.Parse(accepted.ProofJson);
                    var revealed_attrs = jToken.SelectToken("requested_proof.revealed_attr_groups.access-requirement.values");
                    var verify = await _proofService.VerifyProofAsync(context, accepted.Id);
                    var req = _requests.Get(model.ClientDID);
                    if (verify && !req.Status.Equals("SC-verification"))
                    {
                        Console.WriteLine("***proof presentation verified {0}", verify);
                        // update the state and the request map
                        req.Status = "SC-verification";
                        _requests.Add(model.ClientDID, req);
                        // invoke api which invokes the chaincode 
                        AppContext.SetSwitch("System.Net.Http.SocketsHttpHandler.Http2UnencryptedSupport", true);
                        /* using var channel = GrpcChannel.ForAddress("http://10.0.0.14:8080"); */
                        using var channel = GrpcChannel.ForAddress("http://10.0.0.1:8081");
                        var client = new Ssidd.SsiddClient(channel);
                        var respMessage = "";
                        switch (req.ClientRequest.Type)
                        {
                            case "read":
                                var request = _readReqInstance(model.ClientDID, revealed_attrs);
                                var reply = await client.ReadAsync(request);

                                await GrpcChannel.ForAddress("http://10.0.0.1:8081").ShutdownAsync();
                                /* await GrpcChannel.ForAddress("http://10.0.0.14:8080").ShutdownAsync(); */
                                if (reply.Accept)
                                {
                                    return Ok(new AccessRequestModel()
                                    {
                                        ConnectionId = model.RequestId,
                                        ResourceId = req.ClientRequest.ResourceId,
                                        Type = req.ClientRequest.Type,
                                        Status = reply.Message,
                                        Data = reply.Data.ToString(Encoding.UTF8),
                                    });
                                }
                                respMessage = reply.Message;
                                break;
                            case "write":
                                var writeRequest = _writeReqIns(model.ClientDID, revealed_attrs);
                                var writeReply = await client.WriteAsync(writeRequest);
                                await GrpcChannel.ForAddress("http://10.0.0.1:8081").ShutdownAsync();
                                /* await GrpcChannel.ForAddress("http://10.0.0.14:8080").ShutdownAsync(); */
                                if (writeReply.Accept)
                                {
                                    return Ok(new AccessRequestModel()
                                    {
                                        ConnectionId = model.RequestId,
                                        ResourceId = req.ClientRequest.ResourceId,
                                        Type = req.ClientRequest.Type,
                                        Status = writeReply.Message,
                                        Data = "",
                                    });
                                }
                                respMessage = writeReply.Message;
                                break;
                            case "delete":
                                Console.WriteLine("**delete access request");
                                break;

                            case "update":
                                Console.WriteLine("**update access request");
                                break;

                            case "Policy":
                                // Policy will also have CRUD operations.
                                Console.WriteLine("**Policy access request");
                                break;

                            case "Organisation":
                                // Organisation will aslo have CRUD operations.
                                Console.WriteLine("**Organisation access request");
                                break;
                            default:
                                Console.WriteLine($"Type of request is not recognized!.");
                                break;

                        }
                        return Ok(new AccessRequestModel()
                        {
                            ConnectionId = model.RequestId,
                            ResourceId = req.ClientRequest.ResourceId,
                            Type = req.ClientRequest.Type,
                            Status = respMessage,
                            Data = "",
                        });
                    }
                }
                return StatusCode(500, new AccessRequestModel()
                {
                    ConnectionId = model.RequestId,
                    Status = "Proof-Rejected"
                });
            }
            catch (Exception e)
            {
                Console.WriteLine("\nException Caught!");
                Console.WriteLine("Message: {0} ", e.Message);
                var st = new StackTrace(e, true);
                var line = st.GetFrame(st.FrameCount - 1).GetFileLineNumber();
                Console.WriteLine("Line number: {0} ", line);
                return null;
            }
        }

        private ReadRequest _readReqInstance(string clientDID, JToken revealed_attrs)
        {
            var req = _requests.Get(clientDID);
            var readRequest = new ReadRequest
            {
                ClientDID = req.ClientDid,
                ResourceID = req.ClientRequest.ResourceId.ToLower(),
                Timestamp = Timestamp.FromDateTime(DateTime.UtcNow),
                IsPolicy = false,
            };
            // TODO: should be dynamically based on the requested attributes in proof presentation                    
            readRequest.ClientAttributes.Add("nftid", Value.ForString(revealed_attrs["nftid"]["raw"].ToString()));
            readRequest.ClientAttributes.Add("organisationid", Value.ForString(revealed_attrs["organisationid"]["raw"].ToString()));
            readRequest.ClientAttributes.Add("expiration", Value.ForString(revealed_attrs["expiration"]["raw"].ToString()));
            if (revealed_attrs["status"]["raw"].ToString().Contains("true"))
            {
                readRequest.ClientAttributes.Add("status", Value.ForBool(true));
            }
            else
            {
                readRequest.ClientAttributes.Add("status", Value.ForBool(false));
            }
            return readRequest;
        }
        private WriteRequest _writeReqIns(string clientDID, JToken revealed_attrs)
        {
            var req = _requests.Get(clientDID);
            var data = ByteString.CopyFrom(req.ClientRequest.Data, Encoding.Unicode);
            var writeRequest = new WriteRequest
            {
                ClientDID = req.ClientDid,
                PolicyID = "p001", // TODO: should be based on fethced policy
                OwnerOrgID = "spain01",
                ResourceID = req.ClientRequest.ResourceId.ToLower(),
                Data = data,
                Timestamp = Timestamp.FromDateTime(DateTime.UtcNow),
                IsPolicy = false,
            };
            // TODO: should be dynamically based on the requested attributes in proof presentation                    
            writeRequest.ClientAttributes.Add("nftid", Value.ForString(revealed_attrs["nftid"]["raw"].ToString()));
            writeRequest.ClientAttributes.Add("organisationid", Value.ForString(revealed_attrs["organisationid"]["raw"].ToString()));
            writeRequest.ClientAttributes.Add("expiration", Value.ForString(revealed_attrs["expiration"]["raw"].ToString()));
            if (revealed_attrs["status"]["raw"].ToString().Contains("true"))
            {
                writeRequest.ClientAttributes.Add("status", Value.ForBool(true));
            }
            else
            {
                writeRequest.ClientAttributes.Add("status", Value.ForBool(false));
            }
            // TODO: should be read from user inputs
            writeRequest.OptionalResAttributes.Add("anykey", Value.ForString("any-value"));
            return writeRequest;
        }


        private PolicyModel FetchPolicy(string type)
        {
            var policyTest = new PolicyModel { };
            switch (type)
            {
                case "read":
                    // TODO: should fetch Policy by invoking a function from PAP chaincode
                    //policyTest = JsonConvert.DeserializeObject<PolicyModel>(CustomPolicy());
                    break;

                case "write":
                    break;

                case "delete":
                    Console.WriteLine("**delete access request");
                    break;
                case "update":
                    Console.WriteLine("**update access request");
                    break;

                case "Policy":
                    // Policy will also have CRUD operations.
                    Console.WriteLine("**Policy access request");
                    break;

                case "Organisation":
                    // Organisation will aslo have CRUD operations.
                    Console.WriteLine("**Organisation access request");
                    break;

                default:
                    Console.WriteLine($"Type of request is not recognized!.");
                    break;
            }
            return policyTest;
        }

        // CustomPolicy: Creates a custom json for Policy for testing purpose (will be retrieved from back-end (PAPSC)) 
        // The policy can also be retrieved from a database connected to this controller.
        private string CustomPolicy()
        {
            var customPolicy = new PolicyModel
            {
                PolicyID = "P-01",
                Rules = new RuleAttributes
                {
                    UserStatus = true,
                    UserExpiration = 1, // 1 DAY
                    UserOrganisationID = "Spain-01",
                    UserNftID = "00231"
                }
            };
            Resource customResource = new Resource
            {
                RescourceID = "r001",
                ResourceAttr = new ResouceAttributes
                {
                    OrganisationID = "Spain-01",
                    NftID = "00231"
                }
            };
            string jsonString = JsonConvert.SerializeObject(customPolicy);
            return jsonString;
        }
    }
}