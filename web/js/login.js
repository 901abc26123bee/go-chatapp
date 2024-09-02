const createUserURL = accountHostName + "/api/account/v1/user"
const loginURL = accountHostName + "/api/account/v1/auth/login"

function showForm(form) {
    document.getElementById('registerFormContainer').classList.remove('active');
    document.getElementById('loginFormContainer').classList.remove('active');
    if (form === 'register') {
        document.getElementById('registerFormContainer').classList.add('active');
    } else {
        document.getElementById('loginFormContainer').classList.add('active');
    }
}

document.getElementById('registerForm').addEventListener('submit', async function(event) {
    event.preventDefault();

    const name = document.getElementById('registerName').value;
    const email = document.getElementById('registerEmail').value;
    const password = document.getElementById('registerPassword').value;

    try {
        const response = await fetch(createUserURL, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ name: name, email: email, password: password })
        });

        if (!response.ok) {
            const errorData = await response.json();
            document.getElementById('registerResponse').textContent = 'Registration failed: ' + (errorData.error || response.statusText);
        } else {
            document.getElementById('registerResponse').textContent = 'Registered successfully!';
        }
    } catch (error) {
        document.getElementById('registerResponse').textContent = 'Registration failed: ' + error.message;
    }
});

document.getElementById('loginForm').addEventListener('submit', async function(event) {
    event.preventDefault();

    const email = document.getElementById('loginEmail').value;
    const password = document.getElementById('loginPassword').value;

    try {
        const response = await fetch(loginURL, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ email: email, password: password })
        });

        const result = await response.json();

        if (!response.ok) {
            document.getElementById('loginResponse').textContent = 'Login failed: ' + (result.error || response.statusText);
        } else {
            document.getElementById('loginResponse').textContent = 'Login successful!';
            // Store the token in sessionStorage
            sessionStorage.setItem('access-token', result.Token);
        }
    } catch (error) {
        document.getElementById('loginResponse').textContent = 'Login failed: ' + error.message;
    }
});
