FROM mcr.microsoft.com/dotnet/sdk:8.0 AS build

WORKDIR /src

COPY ./performance.csproj .
RUN dotnet restore "./performance.csproj"

COPY . .
RUN dotnet build "performance.csproj" -c Release -o /build
RUN dotnet publish "performance.csproj" -c Release -o /publish

FROM mcr.microsoft.com/dotnet/aspnet:8.0

WORKDIR /app
COPY --from=build /publish .

ENTRYPOINT ["dotnet", "performance.dll"]
EXPOSE 8085