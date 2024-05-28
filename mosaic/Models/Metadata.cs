using System.Collections.Generic;

namespace Voltaserve.Mosaic.Models
{
    public class Metadata
    {
        public int Width { get; set; }

        public int Height { get; set; }

        public string Extension { get; set; }

        public List<ZoomLevel> ZoomLevels { get; set; }
    }
}