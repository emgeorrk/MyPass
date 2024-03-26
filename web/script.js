function login() {
    const username = document.querySelector('input[name="username"]').value;
    const password = document.querySelector('input[name="password"]').value;

    const credentials = username + ':' + password;
    const encodedCredentials = btoa(credentials);

    let xhr = new XMLHttpRequest();

    xhr.open('GET', 'http://localhost:56821', true);
    xhr.setRequestHeader('Authorization', 'Basic ' + encodedCredentials);
    xhr.setRequestHeader('Action', 'getBase')

    let promise = new Promise(function (resolve, reject) {
        xhr.onload = function () {
            if (xhr.status >= 200 && xhr.status < 300) {
                resolve(xhr.responseText);
            } else {
                reject(xhr.status);
            }
        };

        xhr.onerror = function () {
            reject(xhr.status);
        };
    });

    xhr.send();

    promise.then(function(response) {
        console.log('Response:', response);
    }).catch(function(error) {
        console.error('Request failed with status:', error);
    });
}