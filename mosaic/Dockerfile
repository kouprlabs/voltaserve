FROM --platform=$BUILDPLATFORM mcr.microsoft.com/dotnet/sdk:8.0 AS build

WORKDIR /src

COPY ./mosaic.csproj .
RUN dotnet restore "./mosaic.csproj"

COPY . .
RUN dotnet build "mosaic.csproj" -c Release -o /build
RUN dotnet publish "mosaic.csproj" -c Release -o /publish

FROM --platform=$BUILDPLATFORM mcr.microsoft.com/dotnet/aspnet:8.0

WORKDIR /app
COPY --from=build /publish .

ENTRYPOINT ["dotnet", "mosaic.dll"]
EXPOSE 8085