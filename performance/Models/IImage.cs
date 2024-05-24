namespace Defyle.Core.Preview.Models
{
	using System.Drawing;
	using System.IO;

	public interface IImage
	{
		int Width { get; }

		int Height { get; }

		string Extension { get; }

		void Load(string file);

		void Crop(int x, int y, int width, int height);

		void Crop(Rectangle rectangle);

		void ScaleWithAspectRatio(int width, int height);

		void Save(string file);

		void SaveAsPngToStream(Stream stream);
	}
}