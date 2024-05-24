namespace Defyle.Core.Preview.Models
{
    public class ZoomLevel
    {
        public int Index { get; set; }

        public int Width { get; set; }

        public int Height { get; set; }

        public int Rows { get; set; }

        public int Cols { get; set; }

        public float ScaleDownPercentage { get; set; }

        public Tile Tile { get; set; }
    }
}