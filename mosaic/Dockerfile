FROM mcr.microsoft.com/dotnet/sdk:8.0 AS build

WORKDIR /src

COPY ./tiling.csproj .
RUN dotnet restore "./tiling.csproj"

COPY . .
RUN dotnet build "tiling.csproj" -c Release -o /build
RUN dotnet publish "tiling.csproj" -c Release -o /publish

FROM mcr.microsoft.com/dotnet/aspnet:8.0

WORKDIR /app
COPY --from=build /publish .

ENTRYPOINT ["dotnet", "tiling.dll"]
EXPOSE 8085