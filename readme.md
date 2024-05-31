# Desafio Open Telemetry

Objetivo: Desenvolver um sistema em Go que receba um CEP, identifica a cidade e retorna o clima atual (temperatura em graus celsius, fahrenheit e kelvin) juntamente com a cidade. Esse sistema deverá implementar OTEL(Open Telemetry) e Zipkin.

Service A: API que recebe o Input do usuário e envia para o Service B
Service B: API que recebe o Input do Service A e que busca as informações necessárias.

# Subir a aplicação localmente/Rodar o Dockerfile

Executar no cmd o seguinte comando na pasta root do projeto: `docker compose up -d`.

# Chamar API de Temperatura via CEP.

Acessar a seguinte cURL `curl --location 'http://localhost:8080/temperature?cep=CEP'` Onde CEP é o CEP desejado.

# Acessar Zipkin para telemetria

Após realizar a(s) requisição(ões), acessar a URL: `http://localhost:9411/zipkin/` e fazer o filtro necessário para rodar a query e analisar a requisição em detalhes.

Se a requisição for bem sucedida, deverá aparecer no retorno do service-b as 3 temperaturas e a cidade do cep.

# (Bônus) Acessar Jaeger UI

Após realizar a(s) requisição(ões), acessar a URL: `http://localhost:16686/` e fazer o filtro necessário para rodar a query e analisar a requisição em detalhes.

Se a requisição for bem sucedida, deverá aparecer no retorno do service-b as 3 temperaturas e a cidade do cep.