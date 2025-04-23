# Comandos necesarios:

## Instalacion de docker en EC2

1. Actualizamos el SO
    ```sh
    sudo apt-get update
    ```

2. Instale los paquetes necesarios para permitir que apt use paquetes a través de HTTPS:
    ```sh
    sudo apt-get install apt-transport-https ca-certificates curl software-properties-common
    ```

3. Agregue la clave GPG oficial de Docker:
    ```sh
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
    ```

4. Agregue el repositorio de Docker:

    ```sh
    echo "deb [signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
    ```

5. Actualizar e instalar docker
    ```sh
    sudo apt-get update
    sudo apt-get install docker-ce docker-ce-cli containerd.io
    ```

## Dockerizar

1. Crear imagen de docker
    ```sh
    docker build -t go-fiber-app .
    ```

2. Ejecutar contenedor
    ```sh
    docker run -p 3000:3000 go-fiber-app
    ```

## Comandos docker

1. Visualizar las imagenes de docker
```sh
docker images
```
2. Visualizar los contenedores activos
```sh
docker ps
```