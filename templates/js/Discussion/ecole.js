const username = "{{ .Username }}";
const socket = new WebSocket("wss://localhost:443/wsEcole");

socket.onopen = (event) => {
    console.log("WebSocket connection opened.");
};
socket.onmessage = (event) => {
    console.log("Received message:", event.data);
    const msg = JSON.parse(event.data);
    displayMessage(msg);
    if (msg.ImageURL) {
        displayImage(msg.ImageURL);
    } else {
        displayMessage(msg);
    }
};

socket.onclose = (event) => {
    console.log("WebSocket connection closed.");
};
function displayMessage(msg) {
    const chatBox = document.getElementById("chatBox");
    const newMessage = document.createElement('div');
    newMessage.classList.add("message");

    const existingMessage = chatBox.querySelector(`[data-id="${msg.ID}"]`);
    if (existingMessage) {
        return;
    }

    const messageContent = document.createElement("p");
    messageContent.innerHTML = `<strong>${msg.Username}:</strong> ${msg.Content}`;
    newMessage.appendChild(messageContent);

    if (msg.Image) {
        const imageElement = document.createElement("img");
        imageElement.src = msg.Image;
        newMessage.appendChild(imageElement);
    }

    const reactionButtons = document.createElement("div");
    reactionButtons.classList.add("reaction-buttons");

    const likeForm = document.createElement("form");
    likeForm.action = "/likesEcole";
    likeForm.method = "post";

    const likeInput = document.createElement("input");
    likeInput.type = "hidden";
    likeInput.name = "id";
    likeInput.value = msg.ID;

    const likeButton = document.createElement("button");
    likeButton.type = "submit";
    likeButton.classList.add("like-button");
    likeButton.textContent = "Like";

    const likesCount = document.createElement("span");
    likesCount.classList.add("likes-count");
    likesCount.textContent = msg.Likes;

    likeForm.appendChild(likeInput);
    likeForm.appendChild(likeButton);
    likeForm.appendChild(likesCount);

    const dislikeForm = document.createElement("form");
    dislikeForm.action = "/dislikesEcole";
    dislikeForm.method = "post";

    const dislikeInput = document.createElement("input");
    dislikeInput.type = "hidden";
    dislikeInput.name = "id";
    dislikeInput.value = msg.ID;

    const dislikeButton = document.createElement("button");
    dislikeButton.type = "submit";
    dislikeButton.classList.add("dislike-button");
    dislikeButton.textContent = "Dislike";

    const dislikesCount = document.createElement("span");
    dislikesCount.classList.add("dislikes-count");
    dislikesCount.textContent = msg.Dislikes;

    dislikeForm.appendChild(dislikeInput);
    dislikeForm.appendChild(dislikeButton);
    dislikeForm.appendChild(dislikesCount);

    reactionButtons.appendChild(likeForm);
    reactionButtons.appendChild(dislikeForm);

    newMessage.appendChild(reactionButtons);

    chatBox.appendChild(newMessage);
}
document.getElementById("sendButton").addEventListener("click", () => {
    const messageInput = document.getElementById("messageInput");
    const message = messageInput.value;
    const likes = 0;
    const dislikes = 0;

    const msg = {
        Username: username,
        Content: message,
        SocketID: socket.id,
        Likes: likes,
        Dislikes: dislikes,
    };

    const fileInput = document.getElementById("inputImage");
    const imageFile = fileInput.files[0];

    if (imageFile) {
        const reader = new FileReader();
        reader.readAsDataURL(imageFile);
        reader.onload = function () {
            msg.Image = reader.result;
            socket.send(JSON.stringify(msg));
            messageInput.value = "";
            displayMessage(msg);
        };

        reader.onerror = function (error) {
            console.error('Erreur de lecture de l\'image:', error);
        };
    } else {
        socket.send(JSON.stringify(msg));
        messageInput.value = "";

        displayMessage(msg);
    }
});
document.getElementById("chatBox").addEventListener("click", async (event) => {
    event.preventDefault();
    const { target } = event;
    if (target.classList.contains('like-button') || target.classList.contains('dislike-button')) {
        const messageDiv = target.closest('.message');
        const likesCount = messageDiv.querySelector('.likes-count');
        const dislikesCount = messageDiv.querySelector('.dislikes-count');
        let url = '';
        if (target.classList.contains('like-button')) {
            likesCount.textContent = parseInt(likesCount.textContent) + 1;
            url = '/likesEcole';
        } else {
            dislikesCount.textContent = parseInt(dislikesCount.textContent) + 1;
            url = '/dislikesEcole';
        }
        const messageID = messageDiv.querySelector('input[name="id"]').value;

        try {
            const response = await fetch(url, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded',
                },
                body: `id=${messageID}`,
            });

            if (!response.ok) {
                throw new Error('Network response was not ok.');
            }
        } catch (error) {
            console.error('Fetch Error:', error);
        }
    }
});
async function fetchLikesDislikes() {
    try {
        const response = await fetch('/LikesDislikesEcole');
        if (!response.ok) {
            throw new Error('Network response was not ok.');
        }
        return await response.json();
    } catch (error) {
        console.error('Fetch Error:', error);
        return [];
    }
}
fetchLikesDislikes().then((likesDislikes) => {
    console.log('Likes and Dislikes:', likesDislikes);
    likesDislikes.forEach((ld) => {
        const messageDiv = document.querySelector(`input[value="${ld.ID}"]`).closest('.message');
        messageDiv.querySelector('.likes-count').textContent = ld.Likes;
        messageDiv.querySelector('.dislikes-count').textContent = ld.Dislikes;
    });
});

function afficherImageSelectionnee(event) {
    const fichier = event.target.files[0];
    const imageElement = document.createElement('img');

    imageElement.onload = function() {
        URL.revokeObjectURL(imageElement.src);
    };
    imageElement.src = URL.createObjectURL(fichier);
    imageElement.style.maxWidth = '150px';

    const imageContainer = document.getElementById('imageContainer');
    imageContainer.innerHTML = '';
    imageContainer.appendChild(imageElement);
}
document.getElementById('inputImage').addEventListener('change', afficherImageSelectionnee);

function envoyerImage() {
    const fichier = document.getElementById('inputImage').files[0];
    if (!fichier) {
        alert('Veuillez sélectionner une image.');
        return;
    }
    const formData = new FormData();
    formData.append('image', fichier);

    fetch('https://localhost:443/uploadEcole', {
        method: 'POST',
        body: formData
    })
        .then(response => {
            if (!response.ok) {
                throw new Error('Erreur lors de l\'envoi de l\'image.');
            }
            return response.text();
        })
        .then(data => {
            alert('Image envoyée avec succès !');
        })
        .catch(error => {
            console.error('Erreur :', error);
        });
}
function displayImage(imageURL) {
    const chatBox = document.getElementById("chatBox");
    const newImage = document.createElement("div");
    newImage.classList.add("message");

    const imageElement = document.createElement("img");
    imageElement.src = imageURL;
    newImage.appendChild(imageElement);

    chatBox.appendChild(newImage);
}
function getImagesForId(id) {
    fetch(`/getImageEcole?id=${id}`)
        .then(response => response.json())
        .then(images => {
            const imageContainer = document.getElementById("imageContainer");
            images.forEach(imageURL => {
                const imageElement = document.createElement("img");
                imageElement.src = imageURL;
                imageElement.alt = "Image";
                imageElement.style.maxWidth = "150px";
                imageContainer.appendChild(imageElement);
            });
        })
        .catch(error => {
            console.error('Erreur lors de la récupération des images :', error);
        });
}
getImagesForId();