using System;
using Hyperledger.Aries.Agents;
using Newtonsoft.Json;

namespace Hyperledger.Aries.Features.PresentProof.Messages
{
    /// <summary>
    /// Presentation acknowledge message
    /// </summary>
    public class PresentationAcknowledgeMessage : AgentMessage
    {
        /// <summary>
        /// Initializes a new instance of the <see cref="PresentationAcknowledgeMessage" /> class.
        /// </summary>
        public PresentationAcknowledgeMessage()
        {
            Id = Guid.NewGuid().ToString();
            Type = UseMessageTypesHttps ? MessageTypesHttps.PresentProofNames.AcknowledgePresentation : MessageTypes.PresentProofNames.AcknowledgePresentation;
        }
        
        /// <summary>
        /// Initializes a new instance of the <see cref="PresentationAcknowledgeMessage" /> class.
        /// </summary>
        public PresentationAcknowledgeMessage(bool useMessageTypesHttps = false) : base(useMessageTypesHttps)
        {
            Id = Guid.NewGuid().ToString();
            Type = UseMessageTypesHttps ? MessageTypesHttps.PresentProofNames.AcknowledgePresentation : MessageTypes.PresentProofNames.AcknowledgePresentation;
        }
        
        /// <summary>
        /// Gets or sets the acknowledgement status.
        /// </summary>
        [JsonProperty("status")]
        public string Status { get; set; }
    }
}
