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
                this.Collections = data.data.Collections || [];
            } catch (error) {
                console.error("Fehler beim Laden der Collections:", error);
            }
        },
        async fetchEntries() {
            try {
                const response = await fetch('/api/v1/entries');
                const data = await response.json();
                this.Entries = data.data.Entries || [];
            } catch (error) {
                console.error("Fehler beim Laden der Entries:", error);
            }
        },
        async deletething(id, name, type) {
            if (!confirm(`Realy delete ${type} ${name}?`)) return;
                try {
                    const response = await fetch(`/api/v1/${type}/${id}`, {
                        method: 'DELETE',
                        headers: { 'Content-Type': 'application/json' }
                    });
                    if (response.ok) {
                        const data = await response.json();
                        console.log("Answer recieved:", data);
                        showSuccessToast(data.message);
                        if (type === "collection") {
                            this.fetchCollections();
                        } 
                        else if (type === "entry") {
                            this.fetchEntries();
                        }
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
createApp({
    delimiters: ['[[', ']]'],
    data() {
        return {
            MediaTypes: {},
            CaseTypes: {},
            Categories: {},
            Genres: {},
            Publishers: {}
        }
    },
    methods: {
        async fetchMediaTypes() {
            try {
                const response = await fetch('/api/v1/mediatypes');
                const data = await response.json();
                this.MediaTypes = data.data.MediaTypes || [];
            } catch (error) {
                console.error("Fehler beim Laden der Media Types:", error);
            }
        },
        async fetchCaseTypes() {
            try {
                const response = await fetch('/api/v1/casetypes');
                const data = await response.json();
                this.CaseTypes = data.data.CaseTypes || [];
            } catch (error) {
                console.error("Fehler beim Laden der Case Types:", error);
            }
        },
        async fetchCategories() {
            try {
                const response = await fetch('/api/v1/categories');
                const data = await response.json();
                this.Categories = data.data.Categories || [];
            } catch (error) {
                console.error("Fehler beim Laden der Categories:", error);
            }
        },
        async fetchGenres() {
            try {
                const response = await fetch('/api/v1/genres');
                const data = await response.json();
                this.Genres = data.data.Genres || [];
            } catch (error) {
                console.error("Fehler beim Laden der Genres:", error);
            }
        },
        async fetchPublisher() {
            try {
                const response = await fetch('/api/v1/publishers');
                const data = await response.json();
                this.Publishers = data.data.Publishers || [];
            } catch (error) {
                console.error("Fehler beim Laden der Publishers:", error);
            }
        },
        async deletething(id, name, type) {
            if (!confirm(`Really delete ${type} ${name}?`)) return;
                try {
                    const response = await fetch(`/api/v1/${type}/${id}`, {
                        method: 'DELETE',
                        headers: { 'Content-Type': 'application/json' }
                    });
                    if (response.ok) {
                        const data = await response.json();
                        console.log("Answer recieved:", data);
                        showSuccessToast(data.message);
                        if (type === "casetype") {
                            this.fetchCaseTypes();
                        } 
                        else if (type === "mediatype") {
                            this.fetchMediaTypes();
                        }
                        else if (type === "category") {
                            this.fetchCategories();
                        }
                        else if (type === "genre") {
                            this.fetchGenres();
                        }
                        else if (type === "publisher") {
                            this.fetchPublishers();
                        }
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
        this.fetchMediaTypes();
        this.fetchCaseTypes();
        this.fetchCategories();
        this.fetchGenres();
        this.fetchPublisher();
    }
}) .mount('#stock');

createApp({
    delimiters: ['[[', ']]'],
    data() {
        return {
            Collection: {
                Category: {}
            }
        }
    },
    methods: {
        async fetchCollection(id) {
            try {
                const response = await fetch(`/api/v1/collection/${id}`);
                const data = await response.json();
                this.Collection = data.data.Collection || [];
            } catch (error) {
                console.error("Fehler beim Laden der Collection:", error);
            }
        }
    },
    mounted() {
        const el = document.querySelector('#collection');
        const id = el.dataset.id;
        this.fetchCollection(id);
    }
}) .mount('#collection');
createApp({
    delimiters: ['[[', ']]'],
    data() {
        return {
            Entry: {
                MediaType: {},
                CaseType: {},
                Collection: {},
                Genre: {},
                Publisher: {}
            }
        }
    },
    methods: {
        async fetchEntry(id) {
            try {
                const response = await fetch(`/api/v1/entry/${id}`);
                const data = await response.json();
                this.Entry = data.data.Entry || [];
            } catch (error) {
                console.error("Fehler beim Laden des Entry:", error);
            }
        }
    },
    mounted() {
        const el = document.querySelector('#entry');
        const id = el.dataset.id;
        this.fetchEntry(id);
    }
}) .mount('#entry');

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