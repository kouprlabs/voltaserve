namespace Voltaserve.Mosaic.Infra
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
