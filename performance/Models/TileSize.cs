namespace Defyle.Core.Preview.Models
{
  using System;
  using System.Drawing;

  public class TileSize
	{
		private Size _value;

		public TileSize(Size value)
		{
			Value = value;
		}

		public Size Value
		{
			get => _value;

			set
			{
				if (!IsValid(value))
				{
					throw new Exception(GetAcceptanceCriteria());
				}

				_value = value;
			}
		}

		public int Width
		{
			get => Value.Width;

			set => Value = new Size(value, Value.Height);
		}

		public int Height
		{
			get => Value.Height;

			set => Value = new Size(Value.Width, value);
		}

		private static bool IsValid(Size value)
		{
			return IsValidWidth(value.Width) && IsValidHeight(value.Height);
		}

		private static bool IsValidWidth(int width)
		{
			return width > 0;
		}

		private static bool IsValidHeight(int height)
		{
			return height > 0;
		}

		private string GetAcceptanceCriteria()
		{
			return $"{nameof(Value.Width)} and {nameof(Value.Height)} of {nameof(TileSize)} should be bigger than zero.";
		}
	}
}