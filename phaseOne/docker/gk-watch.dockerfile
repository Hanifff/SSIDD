FROM streetcred/dotnet-indy:1.14.2

ENV DOTNET_USE_POLLING_FILE_WATCHER 1
WORKDIR /app

COPY . .

WORKDIR /app/agents/gk-pep
RUN dotnet restore "WebAgent.csproj" \
    -s "https://api.nuget.org/v3/index.json"


RUN dotnet build "WebAgent.csproj" -c Release


ENTRYPOINT dotnet watch run  --urls=http://+:7020 --project WebAgent.csproj
