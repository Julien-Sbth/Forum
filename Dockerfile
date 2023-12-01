# image de base Golang
FROM golang:latest

# répertoire de travail dans l'image
WORKDIR /Forum

# Copie des fichiers du projet dans l'image
COPY . .

# Copie des fichiers de certificat dans l'image
COPY KeyHTTPS /Forum/KeyHTTPS

# Exposition du port 443
EXPOSE 443

# Exécution de la commande de construction du projet
RUN go build -o Forum

# commande de démarrage de l'application
CMD ["go", "run", "main.go"]
