namespace Defyle.Core.Preview.Models
{
    using System;
    using System.Drawing;

    public class MinimumScaleSize
    {
        public MinimumScaleSize(Size value)
        {
            if (!IsValid(value))
            {
                throw new Exception(GetAcceptanceCriteria());
            }
            Value = value;
        }

        public Size Value { get; }

        public int Width => Value.Width;

        public int Height => Value.Height;

        private static bool IsValid(Size value)
        {
            return value.Width > 0 && value.Height > 0;
        }

        private string GetAcceptanceCriteria()
        {
            return $"{nameof(Value.Width)} and {nameof(Value.Height)} of {nameof(MinimumScaleSize)} should be bigger than zero.";
        }
    }
}