namespace Defyle.Core.Util
{
  using System;

  public static class DateTimeUtil
  {
    public static DateTime UnixTimeStampToDateTime(double unixTimeStamp)
    {
      DateTime dateTime = new DateTime(1970, 1, 1, 0, 0, 0, 0, System.DateTimeKind.Utc);
      dateTime = dateTime.AddSeconds(unixTimeStamp);
      return dateTime;
    }
  }
}