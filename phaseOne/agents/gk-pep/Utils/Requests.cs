using System;
using System.Collections.Generic;
using WebAgent.Models;
using System.Diagnostics;


namespace WebAgent.Utils
{
    public class Requests
    {
        private Dictionary<string, RequestorModel> _requests;
        public Requests()
        {
            _requests = new Dictionary<string, RequestorModel>();
        }

        public int GetLength()
        {
            return _requests.Count;
        }

        public Dictionary<string, RequestorModel> GetRequests()
        {
            return _requests;
        }

        public RequestorModel Get(string key)
        {
            try
            {
                return _requests[key];
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
        public void Add(string key, RequestorModel req)
        {
            _requests[key] = req;
        }
        public void Remove(string key)
        {
            _requests.Remove(key);
        }
    }

}