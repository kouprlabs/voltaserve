namespace Voltaserve.Mosaic.Models
{
    using System.Drawing;
    using System.IO;
    using ImageMagick;

    public class Image
    {
        private readonly MagickImage _magickImage;
        private readonly string _file;

        public Image(string file)
        {
            _magickImage = new MagickImage();
            _file = file;
            Load(file);
        }

        public Image(Image source)
        {
            if (source.GetType() == typeof(Image))
            {
                if (source is Image concreteSource)
                {
                    _magickImage = new MagickImage(concreteSource._magickImage);
                }
            }
        }

        public int Width => _magickImage.Width;

        public int Height => _magickImage.Height;

        public string Extension => Path.GetExtension(_file);

        public void Load(string file)
        {
            _magickImage.Read(file);
            _magickImage.AutoOrient();
        }

        public void Crop(int x, int y, int width, int height)
        {
            _magickImage.Crop(new MagickGeometry(x, y, width, height));
        }

        public void Crop(Rectangle rectangle)
        {
            Crop(rectangle.X, rectangle.Y, rectangle.Width, rectangle.Height);
        }

        public void ScaleWithAspectRatio(int width, int height)
        {
            var geometry = new MagickGeometry
            {
                IgnoreAspectRatio = false
            };
            if (width != 0)
            {
                geometry.Width = width;
            }
            if (height != 0)
            {
                geometry.Height = height;
            }
            _magickImage.Scale(geometry);
        }

        public void Save(string file)
        {
            _magickImage.Write(file);
        }

        public void SaveAsPngToStream(Stream stream)
        {
            _magickImage.Write(stream);
        }
    }
}