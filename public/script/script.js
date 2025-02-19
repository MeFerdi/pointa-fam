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
function navigateToHomepage(event) {
    event.preventDefault();
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

 // Toggle Add Product Form
 function toggleAddProductForm() {
    const form = document.getElementById('add-product-form');
    form.classList.toggle('hidden');
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
    const headers = token ? { 'Authorization': `Bearer ${token}` } : {};
    const response = await fetch(`/api/products?category=${category}`, { headers });

    if (response.status === 401) {
        alert('Unauthorized. Please log in again.');
        window.location.href = '/login';
        return;
    }

    const products = await response.json();
    const productsList = document.getElementById('products-list');
    productsList.innerHTML = products.map(product => `
        <div class="product-card bg-white p-4 rounded-lg shadow-md">
            <img src="${product.imageURL}" alt="${product.name}" class="w-full h-48 object-cover rounded-t-lg">
            <div class="p-4">
                <h3 class="text-xl font-bold mb-2">${product.name}</h3>
                <p class="text-gray-600">${product.description}</p>
                <p class="text-gray-800 font-bold">$${product.price}</p>
                <p class="text-gray-600">${product.quantity} in stock</p>
                <span class="inline-block bg-green-200 text-green-800 text-xs px-2 py-1 rounded-full">${product.category}</span>
                ${token ? `<button onclick="addToCart(${product.id})" class="bg-green-500 text-white p-2 rounded hover:bg-green-600 transition duration-200 mt-2">Add to Cart</button>` : ''}
            </div>
        </div>
    `).join('');
}

// Call loadProducts when the page loads
document.addEventListener('DOMContentLoaded', () => loadProducts());

// Load Cart
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
    cart = cartItems.map(item => ({
        id: item.id,
        name: item.product.name,
        price: item.product.price,
        quantity: item.quantity
    }));
    updateCartDisplay();
}

// Add to Cart
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

// Remove from Cart
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

// Toggle Cart Page
function toggleCart() {
    const cartPage = document.querySelector('.cart-page');
    const overlay = document.querySelector('.overlay');
    cartPage.classList.toggle('open');
    overlay.classList.toggle('active');
    updateCartDisplay();
}

// Update Cart Display
function updateCartDisplay() {
    const cartItemsContainer = document.getElementById('cart-items');
    const cartTotal = document.getElementById('cart-total');
    const cartCount = document.getElementById('cart-count');
    const cartEmptyMessage = document.getElementById('cart-empty');

    cartItemsContainer.innerHTML = '';
    let total = 0;

    cart.forEach((item, index) => {
        total += item.price * item.quantity;
        const cartItem = document.createElement('div');
        cartItem.className = 'cart-item flex items-center justify-between';
        cartItem.innerHTML = `
            <div>
                <h3 class="text-lg font-semibold">${item.name}</h3>
                <p class="text-gray-600">$${item.price.toFixed(2)}/kg</p>
            </div>
            <div class="flex items-center">
                <button class="bg-red-500 text-white px-2 py-1 rounded-lg hover:bg-red-600 transition duration-200" onclick="removeFromCart(${item.id})">Remove</button>
                <input type="number" value="${item.quantity}" min="1" class="w-16 ml-4 p-1 border border-gray-300 rounded" onchange="updateQuantity(${index}, this)">
            </div>
        `;
        cartItemsContainer.appendChild(cartItem);
    });

    cartTotal.textContent = total.toFixed(2);
    cartCount.textContent = cart.length;

    if (cart.length === 0) {
        cartEmptyMessage.style.display = 'block';
    } else {
        cartEmptyMessage.style.display = 'none';
    }
}

// Update Quantity
function updateQuantity(index, input) {
    const item = cart[index];
    item.quantity = parseInt(input.value);
    updateCartDisplay();
}

// Initialize Cart Display
document.addEventListener('DOMContentLoaded', loadCart);

// Handle Form Submission
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
        alert('Form submitted successfully');
        form.reset();
    } else {
        const error = await response.json();
        alert(`Error: ${error.error}`);
    }
}

// Logout function
function logout() {
    localStorage.removeItem('userID');
    localStorage.removeItem('token');
    localStorage.removeItem('role');
    window.location.href = '/';
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
        toggleAddProductForm();
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