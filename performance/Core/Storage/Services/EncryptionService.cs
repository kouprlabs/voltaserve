namespace Defyle.Core.Storage.Services
{
  using System;
  using System.IO;
  using System.Security.Cryptography;
  using System.Text;
  using Workspace.Models;

  public class EncryptionService
	{
		private const int KeySize = 256;
		private const int BufferSize = 1048576; // 1MB

		public static void EncryptStreamWithSalt(Stream inputStream, Stream outputStream, byte[] password)
		{
			byte[] salt = new byte[16];
			RNGCryptoServiceProvider provider = new RNGCryptoServiceProvider();
			provider.GetBytes(salt);

			var key = new Rfc2898DeriveBytes(password, salt, 52768);

      using var aes = new AesManaged();
      Array.Clear(password, 0, password.Length);
      Array.Resize(ref password, 1);

      // ReSharper disable once RedundantAssignment
      password = null;

      aes.KeySize = KeySize;
      aes.Key = key.GetBytes(aes.KeySize / 8);
      aes.IV = key.GetBytes(aes.BlockSize / 8);
      aes.Padding = PaddingMode.ISO10126;
      aes.Mode = CipherMode.CBC;

      outputStream.Write(salt, 0, salt.Length);

      using CryptoStream cs = new CryptoStream(outputStream, aes.CreateEncryptor(), CryptoStreamMode.Write);
      key.Dispose();

      byte[] buffer = new byte[BufferSize];
      int read;

      while ((read = inputStream.Read(buffer, 0, buffer.Length)) > 0)
      {
        cs.Write(buffer, 0, read);
      }
    }

		public static void DecryptStreamWithSalt(Stream inputStream, Stream outputStream, byte[] password)
		{
			byte[] salt = new byte[16];

			inputStream.Read(salt, 0, salt.Length);

			var key = new Rfc2898DeriveBytes(password, salt, 52768);

      using var aes = new AesManaged();
      Array.Clear(password, 0, password.Length);
      Array.Resize(ref password, 1);

      // ReSharper disable once RedundantAssignment
      password = null;

      aes.KeySize = KeySize;
      aes.Key = key.GetBytes(aes.KeySize / 8);
      aes.IV = key.GetBytes(aes.BlockSize / 8);
      aes.Padding = PaddingMode.ISO10126;
      aes.Mode = CipherMode.CBC;

      using CryptoStream cs = new CryptoStream(inputStream, aes.CreateDecryptor(), CryptoStreamMode.Read);
      key.Dispose();

      byte[] buffer = new byte[BufferSize];

      while (cs.Read(buffer, 0, buffer.Length) > 0)
      {
        outputStream.Write(buffer, 0, buffer.Length);
      }
    }

		public static byte[] EncryptBytesWithSalt(byte[] clear, byte[] password, byte[] salt)
		{
			var key = new Rfc2898DeriveBytes(password, salt, 52768);

      using var aes = new AesManaged();
      Array.Clear(password, 0, password.Length);
      Array.Resize(ref password, 1);

      // ReSharper disable once RedundantAssignment
      password = null;

      aes.KeySize = KeySize;
      aes.Key = key.GetBytes(aes.KeySize / 8);
      aes.IV = key.GetBytes(aes.BlockSize / 8);
      aes.Padding = PaddingMode.PKCS7;
      aes.Mode = CipherMode.CBC;

      byte[] encrypted;

      using (MemoryStream ms = new MemoryStream())
      {
        using (CryptoStream cs = new CryptoStream(ms, aes.CreateEncryptor(), CryptoStreamMode.Write))
        {
          key.Dispose();

          cs.Write(clear, 0, clear.Length);
        }

        encrypted = ms.ToArray();
      }

      return encrypted;
    }

		public static byte[] DecryptBytesWithSalt(byte[] cipherText, byte[] password, byte[] salt)
    {
      using var key = new Rfc2898DeriveBytes(password, salt, 52768);
      using var aes = new AesManaged();
      Array.Clear(password, 0, password.Length);
      Array.Resize(ref password, 1);

      // ReSharper disable once RedundantAssignment
      password = null;

      aes.KeySize = KeySize;
      aes.Key = key.GetBytes(aes.KeySize / 8);
      aes.IV = key.GetBytes(aes.BlockSize / 8);
      aes.Padding = PaddingMode.PKCS7;
      aes.Mode = CipherMode.CBC;

      byte[] decrypted;

      using (MemoryStream ms = new MemoryStream())
      {
        using (CryptoStream cs = new CryptoStream(ms, aes.CreateDecryptor(), CryptoStreamMode.Write))
        {
          cs.Write(cipherText, 0, cipherText.Length);
        }

        decrypted = ms.ToArray();
      }

      return decrypted;
    }

		public static string DecryptBytes(byte[] cipherText, byte[] key, byte[] iv)
		{
			string plaintext;

			using (var aes = new AesManaged())
			{
				aes.Mode = CipherMode.CBC;
				aes.Padding = PaddingMode.PKCS7;
				aes.KeySize = 128;
				aes.Key = key;
				aes.IV = iv;

				var decryptor = aes.CreateDecryptor();

        using MemoryStream ms = new MemoryStream(cipherText);
        using CryptoStream cs = new CryptoStream(ms, decryptor, CryptoStreamMode.Read);
        using var sr = new StreamReader(cs);
        plaintext = sr.ReadToEnd();
      }

			return plaintext;
		}

		public static byte[] EncryptionKey(Workspace workspace, string cipherPassword)
		{
			string password = DecryptBytes(
				Convert.FromBase64String(cipherPassword),
				Convert.FromBase64String(workspace.TransitKey),
				Convert.FromBase64String(workspace.TransitIv));

			return DecryptBytesWithSalt(
				Convert.FromBase64String(workspace.CipherKey),
				Encoding.UTF8.GetBytes(password),
				Convert.FromBase64String(workspace.Salt));
		}
	}
}