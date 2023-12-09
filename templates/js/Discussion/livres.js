function getCookie(name) {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop().split(';').shift();
}

const username = getCookie("username") || "Default Username";

const socket = new WebSocket("wss://localhost:443/wsLivres");
const ID = "chatBox";

socket.onopen = (event) => {
    console.log("WebSocket connection opened.");
};
socket.onmessage = (event) => {
    console.log("Received message:", event.data);
    const msg = JSON.parse(event.data);
    if (msg.ImageURL) {
        displayImage(msg.ImageURL);
    } else {
        displayReceivedMessage(msg);
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
    likeForm.action = "/likesLivres";
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
    dislikeForm.action = "/dislikesLivres";
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
            url = '/likesLivres';
        } else {
            dislikesCount.textContent = parseInt(dislikesCount.textContent) + 1;
            url = '/dislikesLivres';
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

            const updatedLikesDislikes = await response.json();
            likesCount.textContent = updatedLikesDislikes.likes;
            dislikesCount.textContent = updatedLikesDislikes.dislikes;
        } catch (error) {
            console.error('Fetch Error:', error);
        }
    }
});

async function fetchLikesDislikes() {
    try {
        const response = await fetch('/LikesDislikesLivres');
        if (!response.ok) {
            throw new Error('Network response was not ok.');
        }
        return await response.json();
    } catch (error) {
        console.error('Fetch Error:', error);
        return [];
    }
}
fetchLikesDislikes().then((likesDislikesArray) => {
    console.log('Likes and Dislikes Array:', likesDislikesArray);

    likesDislikesArray.forEach((likesDislikes) => {
        const { Likes, Dislikes, ID } = likesDislikes;

        const messageDiv = document.querySelector(`[data-id="${ID}"]`);
        if (messageDiv) {
            messageDiv.querySelector('.likes-count').textContent = Likes;
            messageDiv.querySelector('.dislikes-count').textContent = Dislikes;
        }
    });
}).catch(error => {
    console.error('Une erreur s\'est produite :', error);
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

    fetch('https://localhost:443/uploadLivres', {
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
            if (data === 'success') {
                alert('Image envoyée avec succès !');
            } else {
                throw new Error('Erreur lors de l\'envoi de l\'image.');
            }
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
    imageElement.src = `data:image/png;base64,${imageURL}`;
    newImage.appendChild(imageElement);

    chatBox.appendChild(newImage);
}


async function getImagesForId(id) {
    try {
        const response = await fetch(`/getImageLivres?id=${id}`);
        if (!response.ok) {
            throw new Error('Network response was not ok.');
        }

        const images = await response.json();

        const imageContainer = document.getElementById("imageContainer");
        imageContainer.innerHTML = '';

        if (images && Array.isArray(images) && images.length > 0) {
            images.forEach(imageURL => {
                const imageElement = document.createElement("img");
                imageElement.src = imageURL;
                imageElement.alt = "Image";
                imageElement.style.maxWidth = "150px";
                imageContainer.appendChild(imageElement);
            });
        } else {
            const noImageMessage = document.createElement("p");
            noImageMessage.textContent = "Aucune image disponible pour cet ID.";
            imageContainer.appendChild(noImageMessage);
        }
    } catch (error) {
        console.error('Erreur lors de la récupération des images :', error);
        const errorMessage = document.createElement("p");
        errorMessage.textContent = "Une erreur s'est produite lors de la récupération des images.";
        const imageContainer = document.getElementById("imageContainer");
        imageContainer.innerHTML = '';
        imageContainer.appendChild(errorMessage);
    }
}
getImagesForId()
    .then(() => {
        console.log('Récupération des images pour tout les ID est terminée.');
    })
    .catch(error => {
        console.error('Une erreur s\'est produite :', error);
    });