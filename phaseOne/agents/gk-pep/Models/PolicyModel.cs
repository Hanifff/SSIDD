using System;


namespace WebAgent
{
    public class PolicyModel
    {
        public string PolicyID { get; set; }
        //public List<PolicyAttributes> Attributes;
        public RuleAttributes Rules;
    }
    public class Subject
    {
        public string SubjectID { get; set; }
        public UserAttributes User;
    }

    public class Resource
    {
        public string RescourceID { get; set; }
        public ResouceAttributes ResourceAttr;
    }

    // NOT USED
    public class PolicyAttributes
    {
        /* public UserAttributes User;
        public ResouceAttributes ResourceAttr; */
        public RuleAttributes Rules;
    }

    public class UserAttributes
    {
        public bool Status { get; set; }
        public DateTime EmploymentStart { get; set; }
        public DateTime Expiration { get; set; }
        public string OrganisationID { get; set; }
        public string NftID { get; set; }
    }
    public class ResouceAttributes
    {
        public string OrganisationID { get; set; }
        public string NftID { get; set; }
        // CRUD operation policy
    }

    public class RuleAttributes
    {

        public bool UserStatus { get; set; }
        public int UserExpiration { get; set; }

        public string UserOrganisationID { get; set; }
        public string UserNftID { get; set; }
    }
}