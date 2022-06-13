FROM streetcred/dotnet-indy:1.14.2

ENV DOTNET_USE_POLLING_FILE_WATCHER 1
WORKDIR /app

COPY . .

WORKDIR /app/agents/ta-agent
RUN dotnet restore "WebAgent.csproj" \
    -s "https://api.nuget.org/v3/index.json"


WORKDIR /app/agents/ta-agent
RUN dotnet build "WebAgent.csproj" -c Release

ENTRYPOINT dotnet watch run  --urls=http://+:7001 --project WebAgent.csproj
