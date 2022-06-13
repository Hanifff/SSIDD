using System;
using System.Text;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Mvc;
using WebAgent.Models;
using System.Net.Http;
using System.Net.Http.Headers;
using Newtonsoft.Json;
using Hyperledger.Aries.Agents;
using Hyperledger.Aries.Configuration;
using Hyperledger.Aries.Extensions;
using Hyperledger.Aries.Features.Handshakes.DidExchange;
using Hyperledger.Aries.Features.IssueCredential;
using Hyperledger.Aries.Models.Records;
using Hyperledger.Aries.Features.Handshakes.Connection;
using Hyperledger.Aries.Storage;
using System.Diagnostics;

namespace WebAgent.Controllers
{

    public class AccessController : Controller
    {
        private readonly IConnectionService _connectionService;
        private readonly IAgentProvider _agentContextProvider;
        private readonly HttpClient _client;
        public AccessController(
            IConnectionService connectionService,
            IAgentProvider agentContextProvider
        )
        {
            _agentContextProvider = agentContextProvider;
            _connectionService = connectionService;
            _client = new HttpClient();
        }


        [HttpGet]

        public async Task<IActionResult> Index()
        {
            var context = await _agentContextProvider.GetContextAsync();
            var connections = await _connectionService.ListAsync(context);
            return View(new AccessFormModel() { Connections = connections });
        }
        [HttpPost]
        public async Task<IActionResult> SubmitRequest()
        {
            try
            {
                _client.DefaultRequestHeaders.Clear();
                _client.DefaultRequestHeaders.Accept.Add(new MediaTypeWithQualityHeaderValue("application/json"));
                AccessRequestModel access = new AccessRequestModel
                {
                    ClientDID = HttpContext.Request.Form["ClientDID"],
                    ResourceId = HttpContext.Request.Form["ResourceId"],
                    ConnectionId = HttpContext.Request.Form["ConnectionId"],
                    Type = HttpContext.Request.Form["Type"],
                    Data = HttpContext.Request.Form["Data"],
                    Policy = HttpContext.Request.Form["Policy"],
                    ResAttributes = HttpContext.Request.Form["Resouce-attributes"],
                    Status = "Initialize-access-request"
                };
                StringContent content = new StringContent(JsonConvert.SerializeObject(access), Encoding.UTF8, "application/json");
                Console.WriteLine("********Sending post request to pep");
                var response = await _client.PostAsync("http://10.0.0.13:7020/Access/AccessRequest/", content);
                if (response != null && response.IsSuccessStatusCode)
                {
                    var responseString = await response.Content.ReadAsStringAsync();
                    var res = JsonConvert.DeserializeObject<AccessRequestModel>(responseString);
                    return View(new AccessRequestModel()
                    {
                        ResourceId = access.ResourceId,
                        Type = access.Type,
                        Status = res.Status
                    });
                }
                return View(new AccessRequestModel()
                {
                    ResourceId = access.ResourceId,
                    Type = access.Type,
                    Status = "Failed"
                });
            }
            catch (HttpRequestException e)
            {
                Console.WriteLine("\nException Caught!");
                Console.WriteLine("Message: {0} ", e.Message);
                var st = new StackTrace(e, true);
                var line = st.GetFrame(st.FrameCount - 1).GetFileLineNumber();
                Console.WriteLine("Line number: {0} ", line);
                return StatusCode(500);
            }
        }
    }
}