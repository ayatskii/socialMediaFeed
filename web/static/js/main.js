// Main JavaScript file for Social Media Feed

document.addEventListener('DOMContentLoaded', function() {
    console.log('Social Media Feed loaded');
    
    // Handle logout form submission
    const logoutForms = document.querySelectorAll('form[action="/logout"]');
    logoutForms.forEach(form => {
        form.addEventListener('submit', function(e) {
            // Clear localStorage on logout
            localStorage.removeItem('token');
            localStorage.removeItem('user');
        });
    });
    
    // Add any client-side functionality here
});

// Like post function
async function likePost(postId) {
    const likeButtons = document.querySelectorAll(`.btn-like[data-post-id="${postId}"]`);
    const dislikeButtons = document.querySelectorAll(`.btn-dislike[data-post-id="${postId}"]`);
    
    // Disable buttons during request
    likeButtons.forEach(btn => btn.disabled = true);
    dislikeButtons.forEach(btn => btn.disabled = true);
    
    try {
        const response = await fetch(`/api/posts/${postId}/like`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            credentials: 'include',
        });
        
        const data = await response.json();
        
        if (response.ok) {
            // Update the like count on the page
            likeButtons.forEach(btn => {
                const countSpan = btn.querySelector('.like-count');
                if (countSpan) {
                    const currentCount = parseInt(countSpan.textContent) || 0;
                    countSpan.textContent = currentCount + 1;
                }
                btn.classList.add('reacted');
                btn.disabled = true;
            });
            
            // If user previously disliked, decrement dislike count and enable dislike button
            dislikeButtons.forEach(btn => {
                const countSpan = btn.querySelector('.dislike-count');
                if (countSpan) {
                    const currentCount = parseInt(countSpan.textContent) || 0;
                    if (currentCount > 0) {
                        countSpan.textContent = currentCount - 1;
                    }
                }
                btn.classList.remove('reacted');
                btn.disabled = false;
            });
        } else {
            if (data.error && data.error.includes('already liked')) {
                // User already liked, disable the button
                likeButtons.forEach(btn => {
                    btn.classList.add('reacted');
                    btn.disabled = true;
                });
            } else if (data.error && data.error.includes('Unauthorized')) {
                alert('Please login to like posts');
            } else {
                alert(data.error || 'Failed to like post');
            }
        }
    } catch (error) {
        console.error('Error liking post:', error);
        alert('An error occurred while liking the post');
        // Re-enable buttons on error
        const isAlreadyLiked = document.querySelector(`.btn-like[data-post-id="${postId}"].reacted`);
        if (!isAlreadyLiked) {
            likeButtons.forEach(btn => btn.disabled = false);
            dislikeButtons.forEach(btn => btn.disabled = false);
        }
    }
}

// Dislike post function
async function dislikePost(postId) {
    const likeButtons = document.querySelectorAll(`.btn-like[data-post-id="${postId}"]`);
    const dislikeButtons = document.querySelectorAll(`.btn-dislike[data-post-id="${postId}"]`);
    
    // Disable buttons during request
    likeButtons.forEach(btn => btn.disabled = true);
    dislikeButtons.forEach(btn => btn.disabled = true);
    
    try {
        const response = await fetch(`/api/posts/${postId}/dislike`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            credentials: 'include',
        });
        
        const data = await response.json();
        
        if (response.ok) {
            // Update the dislike count on the page
            dislikeButtons.forEach(btn => {
                const countSpan = btn.querySelector('.dislike-count');
                if (countSpan) {
                    const currentCount = parseInt(countSpan.textContent) || 0;
                    countSpan.textContent = currentCount + 1;
                }
                btn.classList.add('reacted');
                btn.disabled = true;
            });
            
            // If user previously liked, decrement like count and enable like button
            likeButtons.forEach(btn => {
                const countSpan = btn.querySelector('.like-count');
                if (countSpan) {
                    const currentCount = parseInt(countSpan.textContent) || 0;
                    if (currentCount > 0) {
                        countSpan.textContent = currentCount - 1;
                    }
                }
                btn.classList.remove('reacted');
                btn.disabled = false;
            });
        } else {
            if (data.error && data.error.includes('already disliked')) {
                // User already disliked, disable the button
                dislikeButtons.forEach(btn => {
                    btn.classList.add('reacted');
                    btn.disabled = true;
                });
            } else if (data.error && data.error.includes('Unauthorized')) {
                alert('Please login to dislike posts');
            } else {
                alert(data.error || 'Failed to dislike post');
            }
        }
    } catch (error) {
        console.error('Error disliking post:', error);
        alert('An error occurred while disliking the post');
        // Re-enable buttons on error
        const isAlreadyDisliked = document.querySelector(`.btn-dislike[data-post-id="${postId}"].reacted`);
        if (!isAlreadyDisliked) {
            likeButtons.forEach(btn => btn.disabled = false);
            dislikeButtons.forEach(btn => btn.disabled = false);
        }
    }
}

