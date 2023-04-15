function showError(error) {
    const errorNode = document.querySelector('#error');
    if (errorNode.firstChild) {
        errorNode.remoteChild(errorNode.firstChild);
    }

    errorNode.appendChild(document.createTextNode(error.message || error));
}

function loadSession() {
    // fetch functionality
}