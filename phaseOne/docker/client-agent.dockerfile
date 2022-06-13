FROM streetcred/dotnet-indy:1.14.2 AS base
WORKDIR /app

FROM streetcred/dotnet-indy:1.14.2 AS build
WORKDIR /src
COPY [".", "."]

WORKDIR /src/agents/client-agent
RUN dotnet restore "WebAgent.csproj" \
    -s "https://api.nuget.org/v3/index.json"

COPY ["agents/client-agent/", "."]
COPY ["docker/docker_pool_genesis.txn", "./pool_genesis.txn"]
RUN dotnet build "WebAgent.csproj" -c Release -o /app

FROM build AS publish
RUN dotnet publish "WebAgent.csproj" -c Release -o /app

FROM base AS final
WORKDIR /app
COPY --from=publish /app .

ENTRYPOINT ["dotnet", "WebAgent.dll"]