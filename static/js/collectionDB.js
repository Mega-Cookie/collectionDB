function showSuccessToast(message) {
    const container = document.getElementById('toast-container');
    const toast = document.createElement('div');
    toast.className = 'SuccessToast';
    toast.innerText = message;
    container.appendChild(toast);
    setTimeout(() => {
        toast.style.opacity = '0';
        toast.style.transition = 'opacity 0.5s ease';
        setTimeout(() => toast.remove(), 500);
    }, 3000);
}
function showErrorToast(message) {
    const container = document.getElementById('toast-container');
    const toast = document.createElement('div');
    toast.className = 'ErrorToast';
    toast.innerText = message;
    container.appendChild(toast);
    setTimeout(() => {
        toast.style.opacity = '0';
        toast.style.transition = 'opacity 0.5s ease';
        setTimeout(() => toast.remove(), 500);
    }, 3000);
}
function getList() {
  fetch('/list')
    .then(response => response.json()) // Wir erwarten JSON
    .then(data => {
            const container = document.getElementById('list')
    })
    .catch(error => console.error('Fehler:', error));
}
document.addEventListener('click', async function(e) {
    if (e.target.classList.contains('delete-entry-btn')) {
        const id = e.target.getAttribute('data-id');
        const name = e.target.getAttribute('data-name')
        if (!confirm(`Realy delete entry ${name}?`)) return;
        try {
            const response = await fetch(`/entries/${id}/delete`, {
                method: 'DELETE',
                headers: { 'Content-Type': 'application/json' }
            });
        if (response.ok) {
            const data = await response.json();
            console.log("Answer recieved:", data);
            const card = e.target.closest('.entry-card');
            if (card) {
                card.remove();
            }
            showSuccessToast(data.message);
            } 
         else {
                const data = await response.json();
                console.log("Answer recieved:", data);
                showErrorToast(data.error);
            }
        } catch (err) {
            console.error("Error:", err);
        }
    }
    if (e.target.classList.contains('delete-collect-btn')) {
        const id = e.target.getAttribute('data-id');
        const name = e.target.getAttribute('data-name')
        if (!confirm(`Realy delete collection ${name}?`)) return;
        try {
            const response = await fetch(`/collections/${id}/delete`, {
                method: 'DELETE',
                headers: { 'Content-Type': 'application/json' }
            });
        if (response.ok) {
            const data = await response.json();
            console.log("Answer recieved:", data);
            const card = e.target.closest('.collection-card');
            if (card) {
                card.remove();
            }
            showSuccessToast(data.message);
            }
         else {
                const data = await response.json();
                console.log("Answer recieved:", data);
                showErrorToast(data.error);
            }
        } catch (err) {
            console.error("Error:", err);
        }
    }
    if (e.target.classList.contains('delete-stock-btn')) {
        const type = e.target.getAttribute('data-type')
        const id = e.target.getAttribute('data-id');
        if (!confirm(`Realy delete ${type} ${id}?`)) return;
        try {
            const response = await fetch(`/${type}/${id}/delete`, {
                method: 'DELETE',
                headers: { 'Content-Type': 'application/json' }
            });
        if (response.ok) {
            const data = await response.json();
            console.log("Answer recieved:", data);
            const card = e.target.closest('.stock-card');
            if (card) {
                card.remove();
            }
            showSuccessToast(data.message); 
            } 
         else {
                const data = await response.json();
                console.log("Answer recieved:", data);
                showErrorToast(data.error);
            }
        } catch (err) {
            console.error("Error:", err);
        }
    }
});
document.addEventListener('DOMContentLoaded', () => {
    console.log('Das HTML ist fertig geladen!');
    getList();
});