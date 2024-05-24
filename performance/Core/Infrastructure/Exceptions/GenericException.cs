namespace Defyle.Core.Infrastructure.Exceptions
{
  using System;
  using System.Collections.Generic;
  using Poco;

  public class GenericException : Exception
  {
    private readonly List<Error> _errors = new List<Error>();

    public GenericException WithError(Error error)
    {
      _errors.Add(error);
      return this;
    }
    
    public GenericException WithErrors(IEnumerable<Error> errors)
    {
      _errors.AddRange(errors);
      return this;
    }

    public IEnumerable<Error> Errors => _errors;
  }
}