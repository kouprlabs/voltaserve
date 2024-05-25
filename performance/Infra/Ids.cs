namespace Defyle.Performance.Infra
{
    using System;

    public static class Ids
    {
        public static string New()
        {
            return Guid.NewGuid().ToString().Replace("-", "");
        }
    }
}
