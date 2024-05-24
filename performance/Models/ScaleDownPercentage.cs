namespace Defyle.Core.Preview.Models
{
  using System;

  public class ScaleDownPercentage
  {
    private double? _factor;

    public ScaleDownPercentage(ushort value)
    {
      if (!IsValid(value))
      {
        throw new Exception(GetAcceptanceCriteria());
      }

      Value = value;
    }

    public ushort Value { get; }

    public double Factor
    {
      get
      {
        if (!_factor.HasValue)
        {
          _factor = Value * 0.01;
        }

        return _factor.Value;
      }
    }

    private static bool IsValid(ushort value)
    {
      return value > 0 && value < 100;
    }

    private static string GetAcceptanceCriteria()
    {
      return $"{nameof(ScaleDownPercentage)} should be exclusively more than 0, and exclusively less then 100.";
    }
  }
}