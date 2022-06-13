using Hyperledger.Aries.Models.Records;
using System.Collections.Generic;
using Hyperledger.Aries.Features.Handshakes.Common;

namespace WebAgent.Models
{
    public class AccessFormModel
    {
        public List<ConnectionRecord> Connections { get; set; }
    }
}