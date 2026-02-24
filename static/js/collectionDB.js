function showSucessToast(message) {
    const container = document.getElementById('toast-container');
    const toast = document.createElement('div');
    toast.className = 'SucessToast';
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

document.addEventListener('click', async function(e) {
    if (e.target.classList.contains('delete-entry-btn')) {
        const id = e.target.getAttribute('data-id');
        if (!confirm(`Eintrag ${id} wirklich löschen?`)) return;
        try {
            const response = await fetch(`/entries/${id}/delete`, {
                method: 'DELETE',
                headers: { 'Content-Type': 'application/json' }
            });
        if (response.ok) {
            const data = await response.json();
            console.log("Antwort erhalten:", data);
            const card = e.target.closest('.entry-card');
            if (card) {
                card.remove();
            }
            showSucessToast(data.message); 
            } 
         else {
                const data = await response.json();
                console.log("Antwort erhalten:", data);
                showErrorToast(data.error);
            }
        } catch (err) {
            console.error("Fehler:", err);
        }
    }
    if (e.target.classList.contains('delete-collect-btn')) {
        const id = e.target.getAttribute('data-id');
        if (!confirm(`Collection ${id} wirklich löschen?`)) return;
        try {
            const response = await fetch(`/collections/${id}/delete`, {
                method: 'DELETE',
                headers: { 'Content-Type': 'application/json' }
            });
        if (response.ok) {
            const data = await response.json();
            console.log("Antwort erhalten:", data);
            const card = e.target.closest('.collection-card');
            if (card) {
                card.remove();
            }
            showSucessToast(data.message); 
            } 
         else {
                const data = await response.json();
                console.log("Antwort erhalten:", data);
                showErrorToast(data.error);
            }
        } catch (err) {
            console.error("Fehler:", err);
        }
    }
    if (e.target.classList.contains('delete-stock-btn')) {
        const type = e.target.getAttribute('data-type')
        const id = e.target.getAttribute('data-id');
        if (!confirm(`${type} ${id} wirklich löschen?`)) return;
        try {
            const response = await fetch(`/${type}/${id}/delete`, {
                method: 'DELETE',
                headers: { 'Content-Type': 'application/json' }
            });
        if (response.ok) {
            const data = await response.json();
            console.log("Antwort erhalten:", data);
            const card = e.target.closest('.stock-card');
            if (card) {
                card.remove();
            }
            showSucessToast(data.message); 
            } 
         else {
                const data = await response.json();
                console.log("Antwort erhalten:", data);
                showErrorToast(data.error);
            }
        } catch (err) {
            console.error("Fehler:", err);
        }
    }
});