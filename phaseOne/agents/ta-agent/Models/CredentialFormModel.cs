using System.Collections.Generic;
using Hyperledger.Aries.Models.Records;
using Hyperledger.Aries.Features.Handshakes.DidExchange;
using Hyperledger.Aries.Features.Handshakes.Common;
namespace WebAgent.Models
{
    public class CredentialFormModel
    {
        public List<DefinitionRecord> CredentialDefinitions { get; set; }
        public List<SchemaRecord> Schemas { get; set; }
        public List<ConnectionRecord> Connections { get; set; }

        public static string DefaultAttributes =
@"[
    { 'name': 'name', 'value': 'Alice Smith' },
    { 'name': 'organisationid', 'value': 'spain01' },
    { 'name': 'status', 'value': 'true' },
    { 'name': 'employmentstart', 'value': '2020-01-01' },
    { 'name': 'expiration', 'value': '2050-01-01' },
    { 'name': 'nftid', 'value': '00231' },
    { 'name': 'timestamp', 'value': '2050-01-01' },
]";
    }
}
