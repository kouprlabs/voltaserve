namespace Defyle.WebApi.Policy.Requests
{
  using System.Collections.Generic;
  using System.ComponentModel.DataAnnotations;

  public class CreateManyPoliciesRequest
  {
    [Required]
    public IEnumerable<CreatePolicyRequest> Objects { get; set; }
  }
}