namespace Voltaserve.Mosaic.Models
{
    public class Region
    {
        public int ColStart { get; set; }

        public int ColEnd { get; set; }

        public int RowStart { get; set; }

        public int RowEnd { get; set; }

        public bool IncludesRemainingTiles { get; set; }

        public bool IsNull()
        {
            if (ColStart == 0 &&
                ColEnd == 0 &&
                RowStart == 0 &&
                RowEnd == 0)
            {
                return true;
            }
            return false;
        }
    }
}