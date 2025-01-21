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
    event.preventDefault(); // Prevent the form from submitting the traditional way

    // Get form data
    const form = event.target;
    const formData = new FormData(form);
    const data = Object.fromEntries(formData.entries());

    try {
        // Send the signup request to the backend
        const response = await fetch('/api/register', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data),
        });

        // Parse the response
        const result = await response.json();

        // Check if the request was successful
        if (response.ok) {
            // Save token and user details in localStorage
            localStorage.setItem('token', result.token);
            localStorage.setItem('userID', result.userID);
            localStorage.setItem('role', result.role);

            // Redirect based on the user's role
            if (result.role === 'farmer') {
                window.location.href = '/farmer/dashboard';
            } else if (result.role === 'retailer') {
                window.location.href = '/retailer/dashboard';
            } else {
                console.error('Unknown role:', result.role);
                alert('Unknown role. Please contact support.');
            }
        } else {
            // Display error message if the request failed
            document.getElementById('signup-result').innerText = result.message || 'Signup failed. Please try again.';
        }
    } catch (error) {
        // Handle network or other errors
        console.error('Error during signup:', error);
        document.getElementById('signup-result').innerText = 'An error occurred. Please try again.';
    }
}

function redirectToDashboard() {
    // Get the user's role from localStorage
    const role = localStorage.getItem('role');

    // Redirect based on the role
    if (role === 'farmer') {
        window.location.href = '/farmer/dashboard';
    } else if (role === 'retailer') {
        window.location.href = '/retailer/dashboard';
    } else {
        // Handle unknown or missing role
        console.error('Unknown or missing role:', role);
        alert('You are not logged in or your role is unknown. Please log in.');
        window.location.href = '/login'; // Redirect to login page
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
    loadMyProducts(); // Ensure this is called to load user-specific products
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
async function loadProductsByCategory(category, containerId) {
    try {
        const response = await fetch(`/api/products/category?category=${category}`);
        if (!response.ok) {
            throw new Error(`Failed to fetch ${category} products`);
        }
        console.log(response)
        const products = await response.json();
        const container = document.getElementById(containerId);
        container.innerHTML = products.map(product => `
            <div class="bg-white shadow-lg rounded-lg overflow-hidden hover:shadow-xl transition duration-200">
                <img src="${product.imageURL}" alt="${product.name}" class="w-full h-60 object-cover">
                <div class="p-4">
                    <h3 class="text-xl font-semibold">${product.name}</h3>
                    <p>${product.description}</p>
                    <p class="text-gray-600">$${product.price}</p>
                    <button onclick="window.location.href='/product/${product.id}'" class="mt-2 bg-green-500 text-white p-2 rounded hover:bg-green-600 transition duration-200">View Details</button>
                </div>
            </div>
        `).join('');
    } catch (error) {
        console.error(error);
        const container = document.getElementById(containerId);
        container.innerHTML = `<p class="text-red-500">Failed to load ${category} products. Please try again later.</p>`;
    }
}
// Load cart items
async function loadCart() {
    try {
        const response = await fetch('/api/cart', {
            headers: {
                'Authorization': `Bearer ${localStorage.getItem('token')}`,
            },
        });

        if (!response.ok) {
            throw new Error('Failed to load cart');
        }

        const cart = await response.json();
        const cartItemsContainer = document.getElementById('cart-items');
        const cartTotal = document.getElementById('cart-total');

        // Clear existing items
        cartItemsContainer.innerHTML = '';

        // Render cart items
        let total = 0;
        cart.items.forEach(item => {
            const itemTotal = item.product.price * item.quantity;
            total += itemTotal;

            cartItemsContainer.innerHTML += `
                <div class="cart-item bg-white p-4 rounded-lg shadow-md">
                    <div class="flex items-center justify-between">
                        <div>
                            <h3 class="text-xl font-bold">${item.product.name}</h3>
                            <p class="text-gray-600">$${item.product.price} x ${item.quantity}</p>
                            <p class="text-gray-800 font-bold">$${itemTotal.toFixed(2)}</p>
                        </div>
                        <div class="flex items-center space-x-4">
                            <input type="number" min="1" value="${item.quantity}" onchange="updateQuantity(${item.id}, this.value)" class="w-16 p-2 border rounded">
                            <button onclick="removeFromCart(${item.id})" class="bg-red-500 text-white p-2 rounded hover:bg-red-600 transition duration-200">
                                <i class="fas fa-trash"></i>
                            </button>
                        </div>
                    </div>
                </div>
            `;
        });

        // Update total
        cartTotal.textContent = total.toFixed(2);
    } catch (error) {
        console.error('Error:', error);
        alert('Failed to load cart. Please try again.');
    }
}

// Update cart item quantity
async function updateQuantity(itemId, quantity) {
    try {
        const response = await fetch(`/api/cart/${itemId}`, {
            method: 'PUT',
            headers: {
                'Authorization': `Bearer ${localStorage.getItem('token')}`,
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ quantity: parseInt(quantity) }),
        });

        if (!response.ok) {
            throw new Error('Failed to update quantity');
        }

        const data = await response.json();
        loadCart(); // Refresh the cart list
    } catch (error) {
        console.error('Error:', error);
        alert('Failed to update quantity. Please try again.');
    }
}

// Remove item from cart
async function removeFromCart(itemId) {
    try {
        const response = await fetch(`/api/cart/${itemId}`, {
            method: 'DELETE',
            headers: {
                'Authorization': `Bearer ${localStorage.getItem('token')}`,
            },
        });

        if (!response.ok) {
            throw new Error('Failed to remove item from cart');
        }

        const data = await response.json();
        loadCart(); // Refresh the cart list
        updateCartCount(); // Update cart count in the header
    } catch (error) {
        console.error('Error:', error);
        alert('Failed to remove item from cart. Please try again.');
    }
}

// Update cart count in the header
async function updateCartCount() {
    try {
        const response = await fetch('/api/cart/count', {
            headers: {
                'Authorization': `Bearer ${localStorage.getItem('token')}`,
            },
        });

        if (!response.ok) {
            throw new Error('Failed to fetch cart count');
        }

        const data = await response.json();
        document.getElementById('cart-count').textContent = data.count;
    } catch (error) {
        console.error('Error:', error);
    }
}

// Show cart modal
function showCart() {
    document.getElementById('cart-modal').classList.remove('hidden');
    loadCart(); // Load cart items when the modal is opened
}

// Hide cart modal
function hideCart() {
    document.getElementById('cart-modal').classList.add('hidden');
}

// Handle checkout
async function checkout() {
    try {
        const response = await fetch('/api/checkout', {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${localStorage.getItem('token')}`,
            },
        });

        if (!response.ok) {
            throw new Error('Failed to checkout');
        }

        const data = await response.json();
        alert('Checkout successful!');
        hideCart(); // Close the cart modal
        updateCartCount(); // Update cart count in the header
    } catch (error) {
        console.error('Error:', error);
        alert('Failed to checkout. Please try again.');
    }
}


async function handleFormSubmit(event, url, method) {
    event.preventDefault();
    const form = event.target;
    const formData = new FormData(form);
    const token = localStorage.getItem('token');
    const userID = localStorage.getItem('userID');

    // Ensure userID is included in the form data
    formData.append('userID', userID);

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
    if (!response.ok) {
        console.error('Failed to fetch user products');
        return;
    }
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