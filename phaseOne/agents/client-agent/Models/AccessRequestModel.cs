namespace WebAgent.Models
{
    public class AccessRequestModel

    {
        public string ClientDID {get; set;}
        public string ResourceId { get; set; }
        public string ConnectionId { get; set; }
        public string Type { get; set; }

        public string Data { get; set; }

        public string Policy { get; set; }
        public string ResAttributes { get; set; }
        public string Status { get; set; }
    }
}