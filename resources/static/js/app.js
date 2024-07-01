function showForm(formId) {
    const forms = document.querySelectorAll('.form');
    forms.forEach(form => form.classList.remove('active'));
    document.getElementById(formId + '-form').classList.add('active');
    document.querySelector('.form-container').style.display = 'block';
}

document.getElementById('register-form').addEventListener('submit', async (e) => {
    e.preventDefault();
    const username = document.getElementById('reg-username').value;
    const password = document.getElementById('reg-password').value;
    const email = document.getElementById('reg-email').value;

    const response = await fetch('/user/register', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ username, password, email })
    });

    const result = await response.json();
    alert(result.message);
});

document.getElementById('login-form').addEventListener('submit', async (e) => {
    e.preventDefault();
    const username = document.getElementById('login-username').value;
    const password = document.getElementById('login-password').value;

    const response = await fetch('/user/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ username, password })
    });

    const result = await response.json();
    alert(result.message);
});

document.getElementById('recover-form').addEventListener('submit', async (e) => {
    e.preventDefault();
    const email = document.getElementById('recover-email').value;

    const response = await fetch('/user/recover', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ email })
    });

    const result = await response.json();
    alert(result.message);
});

document.getElementById('verify-form').addEventListener('submit', async (e) => {
    e.preventDefault();
    const email = document.getElementById('verify-email').value;
    const code = document.getElementById('verify-code').value;

    const response = await fetch('/user/verify', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ email, code })
    });

    const result = await response.json();
    alert(result.message);
});
