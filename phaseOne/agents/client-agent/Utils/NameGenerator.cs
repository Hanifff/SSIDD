using System;
using System.Globalization;

namespace WebAgent.Utils
{
    public class NameGenerator
    {
        // Shamelessly stolen from https://github.com/moby/moby/blob/master/pkg/namesgenerator/names-generator.go

        public static string GetRandomName()
        {
            var random = new Random();
            return CultureInfo.InvariantCulture.TextInfo.ToTitleCase(
                $"{_left[random.Next(_left.Length)]} " +
                $"{_right[random.Next(_right.Length)]}");
        }

        private static string[] _left =
        {
            "Alice"
        };

        private static string[] _right =
        {
            "Smith"
        };
    }
}