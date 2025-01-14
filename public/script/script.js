document.addEventListener('DOMContentLoaded', () => {
    const userID = localStorage.getItem('userID');
    const token = localStorage.getItem('token');
    const role = localStorage.getItem('role');

    if (userID && token) {
        document.getElementById('user-account').innerHTML = `
            <button onclick="redirectToDashboard()" class="bg-blue-500 text-white p-2 rounded hover:bg-blue-600 transition duration-200">Dashboard</button>
            <button onclick="logout()" class="bg-red-500 text-white p-2 rounded hover:bg-red-600 transition duration-200 ml-2">Logout</button>
        `;
    }

    const messages = ["Welcome to PointaFam", "Connecting farmers with consumers"];
    let currentMessageIndex = 0;
    let currentCharIndex = 0;
    let isDeleting = false;
    const typingSpeed = 150;
    const deletingSpeed = 100;
    const typingElement = document.getElementById('welcome-message');

    function type() {
        if (currentCharIndex < messages[currentMessageIndex].length) {
            typingElement.textContent += messages[currentMessageIndex].charAt(currentCharIndex);
            currentCharIndex++;
            setTimeout(type, typingSpeed);
        } else {
            isDeleting = true;
            setTimeout(deleteText, 1000);
        }
    }

    function deleteText() {
        if (currentCharIndex > 0) {
            typingElement.textContent = messages[currentMessageIndex].substring(0, currentCharIndex - 1);
            currentCharIndex--;
            setTimeout(deleteText, deletingSpeed);
        } else {
            isDeleting = false;
            currentMessageIndex = (currentMessageIndex + 1) % messages.length;
            setTimeout(type, 500);
        }
    }

    type();
});

async function handleSignUp(event) {
    event.preventDefault();
    const form = event.target;
    const formData = new FormData(form);
    const data = Object.fromEntries(formData.entries());

    const response = await fetch('/api/register', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
    });

    const result = await response.json();
    if (response.ok) {
        localStorage.setItem('token', result.token);
        localStorage.setItem('userID', result.userID);
        localStorage.setItem('role', result.role);
        window.location.href = result.role === 'farmer' ? '/farmer/dashboard' : '/retailer/dashboard';
    } else {
        document.getElementById('signup-result').innerText = result.message;
    }
}

function redirectToDashboard() {
    const role = localStorage.getItem('role');
    if (role === 'farmer') {
        window.location.href = '/farmer/dashboard';
    } else if (role === 'retailer') {
        window.location.href = '/retailer/dashboard';
    }
}

function logout() {
    localStorage.removeItem('userID');
    localStorage.removeItem('token');
    localStorage.removeItem('role');
    window.location.href = '/';
}

function subscribe() {
    const emailInput = document.getElementById('email-input');
    const email = emailInput.value;
    const messageElement = document.getElementById('subscription-message');

    if (email) {
        fetch('/subscribe', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ email: email })
        })
        .then(response => response.json())
        .then(data => {
            messageElement.textContent = data.success ? 'Thank you for subscribing!' : 'Subscription failed. Please try again.';
            if (data.success) emailInput.value = '';
        })
        .catch(() => {
            messageElement.textContent = 'An error occurred. Please try again.';
        });
    } else {
        messageElement.textContent = 'Please enter a valid email address.';
    }
}

document.addEventListener('DOMContentLoaded', async () => {
    await loadUserProfile();
    loadProducts();
    loadCart();
});

async function loadUserProfile() {
    const userID = localStorage.getItem('userID');
    const token = localStorage.getItem('token');
    if (userID && token) {
        const response = await fetch(`/api/user/${userID}`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        if (response.ok) {
            const user = await response.json();
            document.getElementById('user-username').innerText = user.username;
            document.getElementById('user-email').innerText = user.email;
            document.getElementById('user-location').innerText = `Location: ${user.location}`;
            document.getElementById('profile-picture').src = user.profilePictureUrl;
            document.getElementById('edit-username').value = user.username;
            document.getElementById('edit-email').value = user.email;
            document.getElementById('edit-location').value = user.location;
        } else {
            console.error('Failed to load user profile');
        }
    } else {
        console.error('User ID or token not found');
    }
}
function toggleDropdown(id) {
    const element = document.getElementById(id);
    if (element.classList.contains('hidden')) {
        element.classList.remove('hidden');
    } else {
        element.classList.add('hidden');
    }
}
async function uploadProfilePicture(event) {
    const file = event.target.files[0];
    if (!file) return;

    const formData = new FormData();
    formData.append('profile_picture', file);

    const userID = localStorage.getItem('userID');
    const token = localStorage.getItem('token');

    const response = await fetch(`/api/user/${userID}/profile-picture`, {
        method: 'POST',
        headers: {
            'Authorization': `Bearer ${token}`
        },
        body: formData
    });

    if (response.ok) {
        const data = await response.json();
        document.getElementById('profile-picture').src = data.profilePictureUrl;
    } else {
        alert('Failed to upload profile picture');
    }
}
async function handleProfileUpdate(event) {
    event.preventDefault();
    const userID = localStorage.getItem('userID');
    const token = localStorage.getItem('token');
    const form = event.target;
    const formData = new FormData(form);

    const response = await fetch(`/api/user/${userID}`, {
        method: 'PUT',
        headers: {
            'Authorization': `Bearer ${token}`
        },
        body: formData
    });

    if (response.ok) {
        alert('Profile updated successfully');
        loadUserProfile();
        toggleDropdown('edit-profile-form');
    } else {
        const error = await response.json();
        alert(`Error: ${error.error}`);
    }
}
async function loadProducts(category = '') {
    const token = localStorage.getItem('token');
    const response = await fetch(`/api/products?category=${category}`, {
        headers: {
            'Authorization': `Bearer ${token}`
        }
    });
    if (response.status === 401) {
        alert('Unauthorized. Please log in again.');
        window.location.href = '/login';
        return;
    }
    const products = await response.json();
    const productsList = document.getElementById('products-list');
    productsList.innerHTML = products.map(product => `
        <div class="bg-white p-4 rounded-lg shadow-md">
            <h3 class="text-xl font-bold mb-2">${product.name}</h3>
            <p class="text-gray-600">${product.description}</p>
            <p class="text-gray-600">$${product.price}</p>
            <p class="text-gray-600">${product.quantity} in stock</p>
            <button onclick="addToCart(${product.id})" class="bg-green-500 text-white p-2 rounded hover:bg-green-600 transition duration-200">Add to Cart</button>
        </div>
    `).join('');
}

async function loadCart() {
    const userID = localStorage.getItem('userID');
    const token = localStorage.getItem('token');
    const response = await fetch(`/api/cart/${userID}`, {
        headers: {
            'Authorization': `Bearer ${token}`
        }
    });
    if (response.status === 500) {
        alert('Error loading cart. Please try again later.');
        return;
    }
    const cartItems = await response.json();
    const cartList = document.getElementById('cart-list');
    cartList.innerHTML = cartItems.length === 0 ? '<p class="text-center text-gray-600">Cart empty</p>' : cartItems.map(item => `
        <div class="bg-white p-4 rounded-lg shadow-md">
            <h3 class="text-xl font-bold mb-2">${item.product.name}</h3>
            <p class="text-gray-600">${item.product.description}</p>
            <p class="text-gray-600">$${item.product.price}</p>
            <p class="text-gray-600">${item.quantity} in cart</p>
            <button onclick="removeFromCart(${item.id})" class="bg-red-500 text-white p-2 rounded hover:bg-red-600 transition duration-200">Remove</button>
        </div>
    `).join('');
}

async function addToCart(productId) {
    const userID = localStorage.getItem('userID');
    const token = localStorage.getItem('token');
    const response = await fetch('/api/cart', {
        method: 'POST',
        headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ productId: productId, userID: parseInt(userID), quantity: 1 })
    });

    if (response.ok) {
        alert('Product added to cart');
        loadCart();
    } else {
        const error = await response.json();
        alert(`Error: ${error.error}`);
    }
}

async function removeFromCart(cartItemId) {
    const token = localStorage.getItem('token');
    const response = await fetch(`/api/cart/${cartItemId}`, {
        method: 'DELETE',
        headers: {
            'Authorization': `Bearer ${token}`
        }
    });

    if (response.ok) {
        alert('Product removed from cart');
        loadCart();
    } else {
        const error = await response.json();
        alert(`Error: ${error.error}`);
    }
}

function showCart() {
    document.getElementById('cart-modal').classList.remove('hidden');
}

function hideCart() {
    document.getElementById('cart-modal').classList.add('hidden');
}

async function handleFormSubmit(event, url, method) {
    event.preventDefault();
    const form = event.target;
    const formData = new FormData(form);
    const token = localStorage.getItem('token');

    const response = await fetch(url, {
        method: method,
        headers: {
            'Authorization': `Bearer ${token}`
        },
        body: formData
    });

    if (response.ok) {
        alert('Product saved successfully');
        loadProducts();
        loadMyProducts();
        form.reset();
        toggleDropdown('add-product-form');
    } else {
        const error = await response.json();
        alert(`Error: ${error.error}`);
    }
}

async function loadMyProducts() {
    const userID = localStorage.getItem('userID');
    const token = localStorage.getItem('token');
    const response = await fetch(`/api/user/${userID}/products`, {
        headers: {
            'Authorization': `Bearer ${token}`
        }
    });
    const products = await response.json();
    const myProductsList = document.getElementById('my-products-list');
    myProductsList.innerHTML = products.map(product => `
        <div class="bg-white p-4 rounded-lg shadow-md">
            <h3 class="text-xl font-bold mb-2">${product.name}</h3>
            <p class="text-gray-600">${product.description}</p>
            <p class="text-gray-600">$${product.price}</p>
            <p class="text-gray-600">${product.quantity} in stock</p>
            <button onclick="populateEditForm(${product.id})" class="bg-blue-500 text-white p-2 rounded hover:bg-blue-600 transition duration-200">Update</button>
            <button onclick="deleteProduct(${product.id})" class="bg-red-500 text-white p-2 rounded hover:bg-red-600 transition duration-200 ml-2">Delete</button>
        </div>
    `).join('');
}

async function populateEditForm(productId) {
    const token = localStorage.getItem('token');
    const response = await fetch(`/api/products/${productId}`, {
        headers: {
            'Authorization': `Bearer ${token}`
        }
    });
    if (response.status === 401) {
        alert('Unauthorized. Please log in again.');
        window.location.href = '/login';
        return;
    }
    const product = await response.json();
    const form = document.getElementById('create-product-form');
    form.querySelector('input[name="id"]').value = product.id;
    form.querySelector('input[name="name"]').value = product.name;
    form.querySelector('input[name="description"]').value = product.description;
    form.querySelector('input[name="price"]').value = product.price;
    form.querySelector('input[name="quantity"]').value = product.quantity;
    form.querySelector('select[name="category"]').value = product.category;
    form.querySelector('input[name="availability"]').checked = product.availability;
    form.querySelector('input[name="farm_id"]').value = product.farm_id;
    toggleDropdown('add-product-form');
}

async function deleteProduct(productId) {
    const token = localStorage.getItem('token');
    const response = await fetch(`/api/products/${productId}`, {
        method: 'DELETE',
        headers: {
            'Authorization': `Bearer ${token}`
        }
    });

    if (response.ok) {
        alert('Product deleted successfully');
        loadProducts();
        loadMyProducts();
    } else {
        const error = await response.json();
        alert(`Error: ${error.error}`);
    }
}