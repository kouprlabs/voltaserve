FROM mcr.microsoft.com/dotnet/core/aspnet:3.1

WORKDIR /app

COPY ./webapi/bin/Release/netcoreapp3.1/linux-x64/publish .
RUN rm -rf ./config

ENTRYPOINT ["dotnet", "./Defyle.WebApi.dll"]

EXPOSE 8080/tcp