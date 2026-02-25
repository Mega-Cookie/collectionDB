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
const { createApp } = Vue;
createApp({
    delimiters: ['[[', ']]'],
    data() {
        return {
            Collections: [],
            Entries: []
        }
    },
    methods: {
        async fetchCollections() {
            try {
                const response = await fetch('/api/v1/collections');
                const data = await response.json();
                this.Collections = data.Collections || [];
            } catch (error) {
                console.error("Fehler beim Laden der Collections:", error);
            }
        },
        async fetchEntries() {
            try {
                const response = await fetch('/api/v1/entries');
                const data = await response.json();
                this.Entries = data.Entries || [];
            } catch (error) {
                console.error("Fehler beim Laden der Entries:", error);
            }
        },
        async deletething(id, name, type) {
            if (!confirm(`Realy delete ${type} ${name}?`)) return;
                try {
                    const response = await fetch(`/${type}/${id}/delete`, {
                        method: 'DELETE',
                        headers: { 'Content-Type': 'application/json' }
                    });
                    if (response.ok) {
                        const data = await response.json();
                        console.log("Answer recieved:", data);
                        showSuccessToast(data.message); 
                        this.fetchCollections();
                        this.fetchEntries();
                    }
                    else {
                            const data = await response.json();
                            console.log("Answer recieved:", data);
                            showErrorToast(data.error);
                        }
                    }
                catch (err) {
                    console.error("Error:", err);
                }
        }
    },
    mounted() {
        this.fetchCollections();
        this.fetchEntries();
    }
}) .mount('#index');
createApp({
    delimiters: ['[[', ']]'],
    data() {
        return {
            Info: []
        }
    },
    methods: {
        async fetchAbout() {
            try {
                const response = await fetch('/api/v1/about');
                const data = await response.json();
                this.Info = data.data.Info || [];
            } catch (error) {
                console.error("Fehler beim Laden der Infos:", error);
            }
        }
    },
    mounted() {
        this.fetchAbout();
    }
}) .mount('#about');
document.addEventListener('click', async function(e) {
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