using System;
using System.Linq;
using System.Threading.Tasks;
using System.Collections.Generic;
using System.Text;
using Newtonsoft.Json;
using Microsoft.AspNetCore.Mvc;
using Hyperledger.Aries.Agents;
using Hyperledger.Aries.Configuration;
using Hyperledger.Aries.Extensions;
using Hyperledger.Aries.Features.Handshakes.Common;
using Hyperledger.Aries.Features.Handshakes.DidExchange;
using Hyperledger.Aries.Features.PresentProof;
using Hyperledger.Aries.Models.Records;
using Hyperledger.Aries.Features.Handshakes.Connection;
using Hyperledger.Aries.Storage;
using WebAgent.Models;
using System.Diagnostics;
using System.Net.Http;
using System.Net.Http.Headers;

namespace WebAgent.Controllers
{
    public class ProofsController : Controller
    {
        private readonly IAgentProvider _agentContextProvider;
        private readonly IProvisioningService _provisionService;
        private readonly IWalletService _walletService;
        private readonly IConnectionService _connectionService;
        private readonly IProofService _proofService;
        private readonly IMessageService _messageService;
        private readonly HttpClient _client;
        public ProofsController(
            IAgentProvider agentContextProvider,
            IProvisioningService provisionService,
            IWalletService walletService,
            IConnectionService connectionService,
            IProofService proofService,
            IMessageService messageService)
        {
            _agentContextProvider = agentContextProvider;
            _provisionService = provisionService;
            _walletService = walletService;
            _connectionService = connectionService;
            _proofService = proofService;
            _messageService = messageService;
            _client = new HttpClient();
        }

        [HttpGet]
        public async Task<IActionResult> Index()
        {
            var context = await _agentContextProvider.GetContextAsync();
            var proofrequests = await _proofService.ListAsync(context);
            var models = new List<ProofViewModel>();
            foreach (var p in proofrequests)
            {
                models.Add(new ProofViewModel
                {
                    ProofId = p.Id.ToString(),
                    ConnectionId = p.ConnectionId,
                    State = p.State.ToString(),
                    RequestJson = p.RequestJson,
                    ProofJson = p.ProofJson
                });
            }
            return View(new ProofsViewModel { Proofs = models });
        }

        [HttpPost]
        public async Task<IActionResult> ProofPresentation(string connectionId, string proofId)
        {
            try
            {
                return View(new ProofPresentationModel { ProofId = proofId, ConnectionId = connectionId });
            }
            catch (Exception e)
            {
                Console.WriteLine("\nException Caught!");
                Console.WriteLine("Message: {0} ", e.Message);
                var st = new StackTrace(e, true);
                var line = st.GetFrame(st.FrameCount - 1).GetFileLineNumber();
                Console.WriteLine("Line number: {0} ", line);
            }
            return RedirectToAction("Index");
        }

        [HttpPost]
        public async Task<IActionResult> SendProof(ProofPresentationModel model)
        {
            try
            {
                // Holder accepts the proof requests and builds a proof.
                var context = await _agentContextProvider.GetContextAsync();
                var connection = await _connectionService.GetAsync(context, model.ConnectionId);
                // Holder stores the proof request
                var proofRecord = await _proofService.GetAsync(context, model.ProofId);
                var proofObject = JsonConvert.DeserializeObject<ProofRequest>(proofRecord.RequestJson);

                var requestedCredentials = new RequestedCredentials();
                foreach (var requestedAttribute in proofObject.RequestedAttributes)
                {
                    var credentials =
                        await _proofService.ListCredentialsForProofRequestAsync(context, proofObject,
                            requestedAttribute.Key);
                    if (credentials == null || !credentials.Any())
                    {
                        return Content("You have no credentails that can proof the access right.");
                    }
                    requestedCredentials.RequestedAttributes.Add(requestedAttribute.Key,
                        new RequestedAttribute
                        {
                            CredentialId = credentials.First().CredentialInfo.Referent,
                            Revealed = true
                        });
                }

                foreach (var requestedAttribute in proofObject.RequestedPredicates)
                {
                    var credentials =
                        await _proofService.ListCredentialsForProofRequestAsync(context, proofObject,
                            requestedAttribute.Key);

                    requestedCredentials.RequestedPredicates.Add(requestedAttribute.Key,
                        new RequestedAttribute
                        {
                            CredentialId = credentials.First().CredentialInfo.Referent,
                            Revealed = true
                        });
                }
                // Holder accepts the proof request and sends a proof
                (var proofMessage, var proofRec) = await _proofService.CreatePresentationAsync(context, model.ProofId, requestedCredentials);
                await _messageService.SendAsync(context, proofMessage, connection);
                var resp = await _NotifyGK(proofRec.Id, model.ConnectionId, connection.MyDid);
                if (resp != null)
                {
                    return View(new AccessRequestModel() // added in the shared dir so it is shared for all controllers
                    {
                        ConnectionId = resp.ConnectionId,
                        ResourceId = resp.ResourceId,
                        Type = resp.Type,
                        Data = resp.Data,
                        Status = resp.Status
                    });
                }
            }
            catch (Exception e)
            {
                Console.WriteLine("\nException Caught!");
                Console.WriteLine("Message: {0} ", e.Message);
                var st = new StackTrace(e, true);
                var line = st.GetFrame(st.FrameCount - 1).GetFileLineNumber();
                Console.WriteLine("Line number: {0} ", line);
                return StatusCode(500);
            }
            return RedirectToAction("Index");
        }

        private async Task<AccessRequestModel> _NotifyGK(string proofId, string connectionId, string clientDID)
        {
            try
            {
                _client.DefaultRequestHeaders.Clear();
                _client.DefaultRequestHeaders.Accept.Add(new MediaTypeWithQualityHeaderValue("application/json"));
                var proofData = new ProofPresentationModel()
                {
                    ProofId = proofId,
                    ConnectionId = connectionId,
                    ClientDID = clientDID
                };
                StringContent content = new StringContent(JsonConvert.SerializeObject(proofData), Encoding.UTF8, "application/json");
                Console.WriteLine("********Sending post request to gk");
                var response = await _client.PostAsync("http://10.0.0.13:7020/Access/ClientProofPresentation/", content);
                var responseString = await response.Content.ReadAsStringAsync();
                var res = JsonConvert.DeserializeObject<AccessRequestModel>(responseString);
                return res;
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
    }
}