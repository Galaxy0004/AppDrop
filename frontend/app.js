/**
 * API service base URL for external communication.
 */
const API_BASE_URL = 'http://localhost:8080';

/**
 * Application-wide state container managing navigation, theme, and resource collections.
 */
const state = {
    pages: [],
    currentPage: 1,
    totalPages: 1,
    perPage: 10,
    selectedPageId: null,
    widgets: [],
    theme: localStorage.getItem('theme') || 'light'
};

/**
 * Registry of primary DOM elements used across the application for direct manipulation.
 */
const elements = {
    sidebar: document.getElementById('sidebar'),
    menuToggle: document.getElementById('menuToggle'),
    themeToggle: document.getElementById('themeToggle'),
    searchInput: document.getElementById('searchInput'),
    addNewBtn: document.getElementById('addNewBtn'),
    sectionTitle: document.getElementById('sectionTitle'),

    dashboardSection: document.getElementById('dashboardSection'),
    pagesSection: document.getElementById('pagesSection'),
    widgetsSection: document.getElementById('widgetsSection'),
    apiDocsSection: document.getElementById('apiDocsSection'),
    previewBtn: document.getElementById('previewBtn'),
    mobileContent: document.getElementById('mobileContent'),
    previewPageTitle: document.getElementById('previewPageTitle'),

    totalPages: document.getElementById('totalPages'),
    totalWidgets: document.getElementById('totalWidgets'),
    homePage: document.getElementById('homePage'),
    apiStatus: document.getElementById('apiStatus'),
    apiStatusBadge: document.getElementById('apiStatusBadge'),
    pageCount: document.getElementById('pageCount'),

    pageModal: document.getElementById('pageModal'),
    widgetModal: document.getElementById('widgetModal'),
    pageModalTitle: document.getElementById('pageModalTitle'),
    widgetModalTitle: document.getElementById('widgetModalTitle'),

    pageForm: document.getElementById('pageForm'),
    widgetForm: document.getElementById('widgetForm'),
    pageId: document.getElementById('pageId'),
    pageName: document.getElementById('pageName'),
    pageRoute: document.getElementById('pageRoute'),
    pageIsHome: document.getElementById('pageIsHome'),
    widgetId: document.getElementById('widgetId'),
    widgetPageId: document.getElementById('widgetPageId'),
    widgetType: document.getElementById('widgetType'),
    widgetPosition: document.getElementById('widgetPosition'),
    widgetConfig: document.getElementById('widgetConfig'),

    recentPagesList: document.getElementById('recentPagesList'),
    pagesGrid: document.getElementById('pagesGrid'),
    pagesPagination: document.getElementById('pagesPagination'),
    widgetsContainer: document.getElementById('widgetsContainer'),
    pageSelector: document.getElementById('pageSelector'),
    widgetTypeFilter: document.getElementById('widgetTypeFilter'),
    pageFilter: document.getElementById('pageFilter'),

    toastContainer: document.getElementById('toastContainer')
};

/**
 * Core initialization routine triggered upon DOM content load.
 */
document.addEventListener('DOMContentLoaded', () => {
    initTheme();
    initEventListeners();
    checkApiStatus();
    loadDashboardData();
});

/**
 * Theme initialization based on persisted user preference.
 */
function initTheme() {
    if (state.theme === 'dark') {
        document.documentElement.setAttribute('data-theme', 'dark');
        updateThemeIcons(true);
    }
}

/**
 * Toggles application theme between light and dark modes and persists the selection.
 */
function toggleTheme() {
    const isDark = state.theme === 'dark';
    state.theme = isDark ? 'light' : 'dark';

    if (isDark) {
        document.documentElement.removeAttribute('data-theme');
    } else {
        document.documentElement.setAttribute('data-theme', 'dark');
    }

    localStorage.setItem('theme', state.theme);
    updateThemeIcons(!isDark);
}

/**
 * Synchronizes UI icons with the current theme state.
 */
function updateThemeIcons(isDark) {
    const moonIcon = elements.themeToggle.querySelector('.icon-moon');
    const sunIcon = elements.themeToggle.querySelector('.icon-sun');

    if (isDark) {
        moonIcon.style.display = 'none';
        sunIcon.style.display = 'block';
    } else {
        moonIcon.style.display = 'block';
        sunIcon.style.display = 'none';
    }
}

/**
 * Registers global event handlers for application interactions.
 */
function initEventListeners() {
    elements.menuToggle.addEventListener('click', () => {
        elements.sidebar.classList.toggle('open');
    });

    elements.themeToggle.addEventListener('click', (e) => {
        e.preventDefault();
        toggleTheme();
    });

    document.querySelectorAll('.nav-item[data-section]').forEach(item => {
        item.addEventListener('click', (e) => {
            e.preventDefault();
            const section = e.currentTarget.dataset.section;
            navigateToSection(section);
        });
    });

    document.querySelectorAll('.view-all[data-section]').forEach(link => {
        link.addEventListener('click', (e) => {
            e.preventDefault();
            const section = e.currentTarget.dataset.section;
            navigateToSection(section);
        });
    });

    elements.addNewBtn.addEventListener('click', handleAddNew);

    elements.previewBtn.addEventListener('click', () => {
        if (state.selectedPageId) {
            scrollToPreview();
        } else {
            showToast('Select a page to preview', 'warning');
        }
    });

    elements.pageSelector.addEventListener('change', (e) => {
        state.selectedPageId = e.target.value;
        if (state.selectedPageId) {
            loadWidgets(state.selectedPageId);
        } else {
            renderEmptyWidgets();
        }
    });

    elements.widgetTypeFilter.addEventListener('change', (e) => {
        if (state.selectedPageId) {
            loadWidgets(state.selectedPageId, e.target.value);
        }
    });

    elements.pageFilter.addEventListener('change', loadPages);

    document.querySelectorAll('.modal-overlay').forEach(overlay => {
        overlay.addEventListener('click', (e) => {
            if (e.target === overlay) {
                overlay.classList.remove('active');
            }
        });
    });

    document.addEventListener('click', (e) => {
        if (window.innerWidth <= 1024 &&
            !elements.sidebar.contains(e.target) &&
            !elements.menuToggle.contains(e.target) &&
            elements.sidebar.classList.contains('open')) {
            elements.sidebar.classList.remove('open');
        }
    });
}


/**
 * Handles application-level navigation between functional sections.
 * @param {string} section - The identifier of the section to display.
 */
function navigateToSection(section) {
    document.querySelectorAll('.nav-item[data-section]').forEach(item => {
        item.classList.remove('active');
        if (item.dataset.section === section) {
            item.classList.add('active');
        }
    });

    document.querySelectorAll('.content-section').forEach(sec => {
        sec.classList.add('hidden');
    });

    switch (section) {
        case 'dashboard':
            elements.sectionTitle.textContent = 'Dashboard';
            elements.dashboardSection.classList.remove('hidden');
            loadDashboardData();
            break;
        case 'pages':
            elements.sectionTitle.textContent = 'Pages';
            elements.pagesSection.classList.remove('hidden');
            loadPages();
            break;
        case 'widgets':
            elements.sectionTitle.textContent = 'Widgets';
            elements.widgetsSection.classList.remove('hidden');
            loadPageSelector();
            break;
        case 'api-docs':
            elements.sectionTitle.textContent = 'API Documentation';
            elements.apiDocsSection.classList.remove('hidden');
            break;
    }

    elements.sidebar.classList.remove('open');
}

/**
 * Orchestrates the creation of new resources based on the currently active application section.
 */
function handleAddNew() {
    const activeSection = document.querySelector('.nav-item.active').dataset.section;

    switch (activeSection) {
        case 'dashboard':
        case 'pages':
            openPageModal();
            break;
        case 'widgets':
            if (state.selectedPageId) {
                openWidgetModal(state.selectedPageId);
            } else {
                showToast('Please select a page first', 'warning');
            }
            break;
        default:
            openPageModal();
    }
}

/**
 * Generic API communication wrapper with error handling and response parsing.
 * @param {string} endpoint - API endpoint path.
 * @param {string} [method='GET'] - HTTP method.
 * @param {Object} [body=null] - Request payload.
 * @returns {Promise<Object>} The parsed JSON response.
 */
async function apiCall(endpoint, method = 'GET', body = null) {
    const options = {
        method,
        headers: {
            'Content-Type': 'application/json'
        }
    };

    if (body) {
        options.body = JSON.stringify(body);
    }

    try {
        const response = await fetch(`${API_BASE_URL}${endpoint}`, options);
        const data = await response.json();

        if (!response.ok) {
            throw { status: response.status, ...data };
        }

        return data;
    } catch (error) {
        if (error.error) {
            throw error;
        }
        throw { error: { code: 'NETWORK_ERROR', message: 'Failed to connect to API' } };
    }
}

/**
 * Validates the operational status of the backend API service.
 * @returns {Promise<boolean>} True if the API is reachable and healthy.
 */
async function checkApiStatus() {
    try {
        const response = await apiCall('/health');
        elements.apiStatus.textContent = 'Connected';
        elements.apiStatusBadge.textContent = 'Online';
        elements.apiStatusBadge.classList.remove('status-inactive');
        elements.apiStatusBadge.classList.add('status-active');
        return true;
    } catch (error) {
        elements.apiStatus.textContent = 'Disconnected';
        elements.apiStatusBadge.textContent = 'Offline';
        elements.apiStatusBadge.classList.remove('status-active');
        elements.apiStatusBadge.classList.add('status-inactive');
        return false;
    }
}

/**
 * Orchestrates the retrieval and display of high-level application metrics and recent activity for the dashboard.
 */
async function loadDashboardData() {
    try {
        const data = await apiCall(`/pages?page=1&per_page=100`);
        state.pages = data.pages || [];

        elements.totalPages.textContent = data.total || 0;
        elements.pageCount.textContent = data.total || 0;

        const homePage = state.pages.find(p => p.is_home);
        elements.homePage.textContent = homePage ? homePage.name : 'Not Set';

        let widgetCount = 0;
        for (const page of state.pages) {
            const pageData = await apiCall(`/pages/${page.id}`);
            widgetCount += (pageData.widgets || []).length;
        }
        elements.totalWidgets.textContent = widgetCount;

        renderRecentPages(state.pages.slice(0, 5));
    } catch (error) {
        console.error('Failed to load dashboard data:', error);
        showToast('Failed to load dashboard data', 'error');
    }
}

/**
 * Renders a succinct list of recently accessed or modified pages.
 * @param {Array} pages - Collection of page objects to render.
 */
function renderRecentPages(pages) {
    if (!pages || pages.length === 0) {
        elements.recentPagesList.innerHTML = `
            <div class="empty-state">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                    <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
                    <polyline points="14 2 14 8 20 8"/>
                </svg>
                <p>No pages created yet</p>
                <button class="btn btn-secondary btn-sm" onclick="openPageModal()">Create First Page</button>
            </div>
        `;
        return;
    }

    elements.recentPagesList.innerHTML = pages.map(page => `
        <div class="recent-item" onclick="viewPage('${page.id}')">
            <div class="recent-item-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
                    <polyline points="14 2 14 8 20 8"/>
                </svg>
            </div>
            <div class="recent-item-content">
                <div class="recent-item-title">${escapeHtml(page.name)}</div>
                <div class="recent-item-meta">${escapeHtml(page.route)}</div>
            </div>
            ${page.is_home ? '<span class="recent-item-badge">Home</span>' : ''}
        </div>
    `).join('');
}


/**
 * Retrieves and renders a paginated collection of pages from the API.
 */
async function loadPages() {
    try {
        const filter = elements.pageFilter.value;
        const data = await apiCall(`/pages?page=${state.currentPage}&per_page=${state.perPage}`);

        let pages = data.pages || [];

        if (filter === 'home') {
            pages = pages.filter(p => p.is_home);
        } else if (filter === 'regular') {
            pages = pages.filter(p => !p.is_home);
        }

        state.pages = pages;
        state.totalPages = data.total_pages || 1;

        renderPages(pages);
        renderPagination();
    } catch (error) {
        console.error('Failed to load pages:', error);
        showToast('Failed to load pages', 'error');
        renderEmptyPages();
    }
}

/**
 * Generates and injects HTML markup for the page cards grid.
 * @param {Array} pages - Collection of page objects to render.
 */
function renderPages(pages) {
    if (!pages || pages.length === 0) {
        renderEmptyPages();
        return;
    }

    elements.pagesGrid.innerHTML = pages.map(page => `
        <div class="page-card ${page.is_home ? 'is-home' : ''}" data-page-id="${page.id}">
            <div class="page-card-header">
                <div class="page-card-title">
                    ${escapeHtml(page.name)}
                    ${page.is_home ? '<span class="home-badge">Home</span>' : ''}
                </div>
                <div class="page-card-actions">
                    <button class="btn-edit" onclick="editPage('${page.id}')" title="Edit">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                            <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
                        </svg>
                    </button>
                    <button class="btn-delete" onclick="deletePage('${page.id}', ${page.is_home})" title="Delete">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <polyline points="3 6 5 6 21 6"/>
                            <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
                        </svg>
                    </button>
                </div>
            </div>
            <div class="page-card-body">
                <div class="page-card-route">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71"/>
                        <path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71"/>
                    </svg>
                    ${escapeHtml(page.route)}
                </div>
                <div class="page-card-stats">
                    <div class="page-card-stat">
                        <span class="page-card-stat-label">Created</span>
                        <span class="page-card-stat-value">${formatDate(page.created_at)}</span>
                    </div>
                    <div class="page-card-stat">
                        <span class="page-card-stat-label">Updated</span>
                        <span class="page-card-stat-value">${formatDate(page.updated_at)}</span>
                    </div>
                </div>
            </div>
            <div class="page-card-footer">
                <button class="btn btn-secondary btn-sm" onclick="viewPageWidgets('${page.id}')">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="14" height="14">
                        <rect x="3" y="3" width="18" height="18" rx="2"/>
                        <line x1="3" y1="9" x2="21" y2="9"/>
                        <line x1="9" y1="21" x2="9" y2="9"/>
                    </svg>
                    View Widgets
                </button>
                <button class="btn btn-primary btn-sm" onclick="openWidgetModal('${page.id}')">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="14" height="14">
                        <line x1="12" y1="5" x2="12" y2="19"/>
                        <line x1="5" y1="12" x2="19" y2="12"/>
                    </svg>
                    Add Widget
                </button>
            </div>
        </div>
    `).join('');
}

/**
 * Renders an empty state placeholder when no pages are available or match the criteria.
 */
function renderEmptyPages() {
    elements.pagesGrid.innerHTML = `
        <div class="empty-state" style="grid-column: 1 / -1;">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
                <polyline points="14 2 14 8 20 8"/>
            </svg>
            <p>No pages found</p>
            <button class="btn btn-primary btn-sm" onclick="openPageModal()">Create Your First Page</button>
        </div>
    `;
}

/**
 * Generates and injects HTML markup for the pagination controls.
 */
function renderPagination() {
    if (state.totalPages <= 1) {
        elements.pagesPagination.innerHTML = '';
        return;
    }

    let html = '';

    html += `<button class="pagination-btn" ${state.currentPage === 1 ? 'disabled' : ''} onclick="goToPage(${state.currentPage - 1})">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="16" height="16">
            <polyline points="15 18 9 12 15 6"/>
        </svg>
    </button>`;

    for (let i = 1; i <= state.totalPages; i++) {
        if (i === 1 || i === state.totalPages || (i >= state.currentPage - 1 && i <= state.currentPage + 1)) {
            html += `<button class="pagination-btn ${i === state.currentPage ? 'active' : ''}" onclick="goToPage(${i})">${i}</button>`;
        } else if (i === state.currentPage - 2 || i === state.currentPage + 2) {
            html += '<span style="padding: 0 8px;">...</span>';
        }
    }

    html += `<button class="pagination-btn" ${state.currentPage === state.totalPages ? 'disabled' : ''} onclick="goToPage(${state.currentPage + 1})">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="16" height="16">
            <polyline points="9 18 15 12 9 6"/>
        </svg>
    </button>`;

    elements.pagesPagination.innerHTML = html;
}

/**
 * Navigates to a specific page within the collection.
 * @param {number} page - The target page number.
 */
function goToPage(page) {
    if (page < 1 || page > state.totalPages) return;
    state.currentPage = page;
    loadPages();
}


/**
 * Populates the page selection dropdown with available pages.
 */
async function loadPageSelector() {
    try {
        const data = await apiCall('/pages?page=1&per_page=100');
        const pages = data.pages || [];

        elements.pageSelector.innerHTML = '<option value="">Select a Page</option>' +
            pages.map(p => `<option value="${p.id}">${escapeHtml(p.name)} (${escapeHtml(p.route)})</option>`).join('');

        if (state.selectedPageId) {
            elements.pageSelector.value = state.selectedPageId;
            loadWidgets(state.selectedPageId);
        }
    } catch (error) {
        console.error('Failed to load page selector:', error);
    }
}

/**
 * Retrieves and renders widgets for a specific page, optionally filtered by type.
 * @param {string} pageId - The identifier of the page.
 * @param {string} [typeFilter=''] - Optional widget type filter.
 */
async function loadWidgets(pageId, typeFilter = '') {
    try {
        let endpoint = `/pages/${pageId}/widgets`;
        if (typeFilter) {
            endpoint += `?type=${typeFilter}`;
        }

        const data = await apiCall(endpoint);
        state.widgets = data.widgets || [];

        renderWidgets(state.widgets);
        renderMobilePreview(state.widgets);
    } catch (error) {
        console.error('Failed to load widgets:', error);
        showToast('Failed to load widgets', 'error');
        renderEmptyWidgets();
    }
}

/**
 * Generates and injects HTML markup for the widget management list.
 * @param {Array} widgets - Collection of widget objects to render.
 */
function renderWidgets(widgets) {
    if (!widgets || widgets.length === 0) {
        elements.widgetsContainer.innerHTML = `
            <div class="empty-state">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                    <rect x="3" y="3" width="18" height="18" rx="2"/>
                    <line x1="3" y1="9" x2="21" y2="9"/>
                </svg>
                <p>No widgets on this page</p>
                <button class="btn btn-primary btn-sm" onclick="openWidgetModal('${state.selectedPageId}')">Add First Widget</button>
            </div>
        `;
        return;
    }

    elements.widgetsContainer.innerHTML = `
        <div class="widgets-list" id="widgetsList">
            ${widgets.map((widget, index) => `
                <div class="widget-item" data-widget-id="${widget.id}" data-type="${widget.type}" draggable="true">
                    <div class="widget-drag-handle">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <circle cx="9" cy="5" r="1"/>
                            <circle cx="9" cy="12" r="1"/>
                            <circle cx="9" cy="19" r="1"/>
                            <circle cx="15" cy="5" r="1"/>
                            <circle cx="15" cy="12" r="1"/>
                            <circle cx="15" cy="19" r="1"/>
                        </svg>
                    </div>
                    <div class="widget-item-icon">
                        ${getWidgetIcon(widget.type)}
                    </div>
                    <div class="widget-item-content">
                        <div class="widget-item-type">${formatWidgetType(widget.type)}</div>
                        <div class="widget-item-position">Position: ${widget.position}</div>
                    </div>
                    <div class="widget-item-actions">
                        <button class="btn-edit" onclick="editWidget('${widget.id}')" title="Edit">
                            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                                <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
                            </svg>
                        </button>
                        <button class="btn-delete" onclick="deleteWidget('${widget.id}')" title="Delete">
                            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                <polyline points="3 6 5 6 21 6"/>
                                <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
                            </svg>
                        </button>
                    </div>
                </div>
            `).join('')}
        </div>
    `;

    initWidgetDragDrop();
}

/**
 * Resets the widget container and preview pane to their default empty states.
 */
function renderEmptyWidgets() {
    elements.widgetsContainer.innerHTML = `
        <div class="empty-state">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <rect x="3" y="3" width="18" height="18" rx="2"/>
                <line x1="3" y1="9" x2="21" y2="9"/>
            </svg>
            <p>Select a page to view its widgets</p>
        </div>
    `;
    elements.mobileContent.innerHTML = `
        <div class="preview-empty">
            <p>Select a page to see the preview</p>
        </div>
    `;
    elements.previewPageTitle.textContent = 'AppDrop';
}


/**
 * Renders a simulated mobile interface displaying current widgets for visual verification.
 * @param {Array} widgets - Collection of widget objects to render in the preview.
 */
function renderMobilePreview(widgets) {
    if (!widgets || widgets.length === 0) {
        elements.mobileContent.innerHTML = `
            <div class="preview-empty">
                <p>No widgets added to this page yet</p>
            </div>
        `;
        return;
    }

    const selectedPage = state.pages.find(p => p.id === state.selectedPageId);
    if (selectedPage) {
        elements.previewPageTitle.textContent = selectedPage.name;
    }

    elements.mobileContent.innerHTML = widgets.map(widget => {
        const config = widget.config || {};

        switch (widget.type) {
            case 'banner':
                return `
                    <div class="preview-banner">
                        ${config.image_url ? `<img src="${config.image_url}" alt="Banner">` : '<div style="height:100%; display:flex; align-items:center; justify-content:center; background:#eee; color:#aaa; font-size:10px;">No Image</div>'}
                        <div class="preview-banner-overlay">
                            <div class="preview-banner-title">${escapeHtml(config.title || 'Welcome Banner')}</div>
                        </div>
                    </div>
                `;
            case 'product_grid':
                const cols = config.columns || 2;
                const limit = config.limit || 4;
                return `
                    <div class="preview-product-grid" style="grid-template-columns: repeat(${cols}, 1fr);">
                        ${Array(limit).fill(0).map(() => `
                            <div class="preview-product-item">
                                <div class="product-image-placeholder"></div>
                                <div class="product-title-placeholder"></div>
                                <div style="height:6px; width:40%; background:#e2e8f0; border-radius:2px; margin-top:4px;"></div>
                            </div>
                        `).join('')}
                    </div>
                `;
            case 'text':
                const style = config.style === 'heading' ? 'font-weight:700; font-size:1rem; margin-bottom:8px;' : '';
                return `
                    <div class="preview-text" style="${style}">
                        ${escapeHtml(config.content || 'Sample Text Content')}
                    </div>
                `;
            case 'image':
                return `
                    <div class="preview-image">
                        ${config.src ? `<img src="${config.src}" alt="Widget Image">` : '<div style="aspect-ratio:16/9; background:#eee; display:flex; align-items:center; justify-content:center; color:#aaa; font-size:10px; border-radius:8px;">No Image</div>'}
                    </div>
                `;
            case 'spacer':
                const height = config.height || 20;
                return `<div class="preview-spacer" style="height: ${height}px;"></div>`;
            default:
                return '';
        }
    }).join('');
}

/**
 * Scrolls the viewport to the mobile preview container on smaller screens.
 */
function scrollToPreview() {
    const previewContainer = document.querySelector('.mobile-preview-container');
    if (window.innerWidth <= 1024) {
        previewContainer.scrollIntoView({ behavior: 'smooth' });
    }
}

/**
 * Returns the SVG markup for a specific widget type.
 * @param {string} type - The widget type identifier.
 * @returns {string} SVG HTML string.
 */
function getWidgetIcon(type) {
    const icons = {
        banner: `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <rect x="3" y="3" width="18" height="18" rx="2"/>
            <circle cx="8.5" cy="8.5" r="1.5"/>
            <polyline points="21 15 16 10 5 21"/>
        </svg>`,
        product_grid: `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <rect x="3" y="3" width="7" height="7"/>
            <rect x="14" y="3" width="7" height="7"/>
            <rect x="3" y="14" width="7" height="7"/>
            <rect x="14" y="14" width="7" height="7"/>
        </svg>`,
        text: `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="4 7 4 4 20 4 20 7"/>
            <line x1="9" y1="20" x2="15" y2="20"/>
            <line x1="12" y1="4" x2="12" y2="20"/>
        </svg>`,
        image: `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <rect x="3" y="3" width="18" height="18" rx="2"/>
            <circle cx="8.5" cy="8.5" r="1.5"/>
            <polyline points="21 15 16 10 5 21"/>
        </svg>`,
        spacer: `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="12" y1="5" x2="12" y2="19"/>
            <polyline points="8 9 12 5 16 9"/>
            <polyline points="8 15 12 19 16 15"/>
        </svg>`
    };
    return icons[type] || icons.text;
}

/**
 * Formats a raw widget type identifier into a human-readable string.
 * @param {string} type - The raw widget type (e.g., 'product_grid').
 * @returns {string} Formatted string (e.g., 'Product Grid').
 */
function formatWidgetType(type) {
    return type.split('_').map(word => word.charAt(0).toUpperCase() + word.slice(1)).join(' ');
}

/**
 * Initializes HTML5 Drag and Drop functionality for reordering widgets.
 */
function initWidgetDragDrop() {
    const list = document.getElementById('widgetsList');
    if (!list) return;

    let draggedItem = null;

    list.querySelectorAll('.widget-item').forEach(item => {
        item.addEventListener('dragstart', (e) => {
            draggedItem = item;
            item.classList.add('dragging');
            e.dataTransfer.effectAllowed = 'move';
        });

        item.addEventListener('dragend', () => {
            item.classList.remove('dragging');
            draggedItem = null;
            saveWidgetOrder();
        });

        item.addEventListener('dragover', (e) => {
            e.preventDefault();
            e.dataTransfer.dropEffect = 'move';

            if (draggedItem && draggedItem !== item) {
                const rect = item.getBoundingClientRect();
                const midY = rect.top + rect.height / 2;

                if (e.clientY < midY) {
                    item.parentNode.insertBefore(draggedItem, item);
                } else {
                    item.parentNode.insertBefore(draggedItem, item.nextSibling);
                }
            }
        });
    });
}


/**
 * Persists the current sequential arrangement of widgets to the persistent store.
 */
async function saveWidgetOrder() {
    if (!state.selectedPageId) return;

    const widgetItems = document.querySelectorAll('#widgetsList .widget-item');
    const widgetIds = Array.from(widgetItems).map(item => item.dataset.widgetId);

    try {
        await apiCall(`/pages/${state.selectedPageId}/widgets/reorder`, 'POST', { widget_ids: widgetIds });
        showToast('Widget order saved', 'success');
        loadWidgets(state.selectedPageId);
    } catch (error) {
        console.error('Failed to save widget order:', error);
        showToast(error.error?.message || 'Failed to save widget order', 'error');
        loadWidgets(state.selectedPageId);
    }
}

/**
 * Interface entry point for creating or editing page resources.
 * @param {Object} [page=null] - Optional page object for edit mode.
 */
function openPageModal(page = null) {
    elements.pageModalTitle.textContent = page ? 'Edit Page' : 'Create Page';
    elements.pageForm.reset();

    if (page) {
        elements.pageId.value = page.id;
        elements.pageName.value = page.name;
        elements.pageRoute.value = page.route;
        elements.pageIsHome.checked = page.is_home;
    } else {
        elements.pageId.value = '';
    }

    elements.pageModal.classList.add('active');
}

/**
 * Dismisses the page management modal.
 */
function closePageModal() {
    elements.pageModal.classList.remove('active');
}

/**
 * Retrieves details for a specific page and opens the edit interface.
 * @param {string} pageId - The identifier of the page to edit.
 */
async function editPage(pageId) {
    try {
        const page = await apiCall(`/pages/${pageId}`);
        openPageModal(page);
    } catch (error) {
        console.error('Failed to load page:', error);
        showToast('Failed to load page details', 'error');
    }
}

/**
 * Validates and persists page resource changes to the API.
 */
async function savePage() {
    const id = elements.pageId.value;
    const name = elements.pageName.value.trim();
    const route = elements.pageRoute.value.trim();
    const isHome = elements.pageIsHome.checked;

    if (!name || !route) {
        showToast('Please fill in all required fields', 'warning');
        return;
    }

    const data = { name, route, is_home: isHome };

    try {
        if (id) {
            await apiCall(`/pages/${id}`, 'PUT', data);
            showToast('Page updated successfully', 'success');
        } else {
            await apiCall('/pages', 'POST', data);
            showToast('Page created successfully', 'success');
        }

        closePageModal();
        loadPages();
        loadDashboardData();
    } catch (error) {
        console.error('Failed to save page:', error);
        showToast(error.error?.message || 'Failed to save page', 'error');
    }
}

/**
 * Facilitates the deletion of a page resource after user confirmation.
 * @param {string} pageId - The identifier of the page to delete.
 * @param {boolean} isHome - Protection flag for the primary home page.
 */
async function deletePage(pageId, isHome) {
    if (isHome) {
        showToast('Cannot delete the home page. Set another page as home first.', 'warning');
        return;
    }

    if (!confirm('Are you sure you want to delete this page? This will also delete all its widgets.')) {
        return;
    }

    try {
        await apiCall(`/pages/${pageId}`, 'DELETE');
        showToast('Page deleted successfully', 'success');
        loadPages();
        loadDashboardData();
    } catch (error) {
        console.error('Failed to delete page:', error);
        showToast(error.error?.message || 'Failed to delete page', 'error');
    }
}

/**
 * Directs the application context to the pages view.
 * @param {string} pageId - The identifier of the page to view.
 */
function viewPage(pageId) {
    navigateToSection('pages');
}

/**
 * Directs the application context to the widgets view for a specific page.
 * @param {string} pageId - The identifier of the page whose widgets will be viewed.
 */
function viewPageWidgets(pageId) {
    state.selectedPageId = pageId;
    navigateToSection('widgets');
}

/**
 * Interface entry point for creating or editing widget resources.
 * @param {string} pageId - The identifier of the parent page.
 * @param {Object} [widget=null] - Optional widget object for edit mode.
 */
function openWidgetModal(pageId, widget = null) {
    elements.widgetModalTitle.textContent = widget ? 'Edit Widget' : 'Add Widget';
    elements.widgetForm.reset();

    elements.widgetPageId.value = pageId;

    if (widget) {
        elements.widgetId.value = widget.id;
        elements.widgetType.value = widget.type;
        elements.widgetPosition.value = widget.position;
        elements.widgetConfig.value = widget.config ? JSON.stringify(widget.config, null, 2) : '';
    } else {
        elements.widgetId.value = '';
    }

    elements.widgetModal.classList.add('active');
}

/**
 * Dismisses the widget management modal.
 */
function closeWidgetModal() {
    elements.widgetModal.classList.remove('active');
}

/**
 * Locates a widget within the local state and opens the edit interface.
 * @param {string} widgetId - The identifier of the widget to edit.
 */
async function editWidget(widgetId) {
    const widget = state.widgets.find(w => w.id === widgetId);
    if (widget) {
        openWidgetModal(state.selectedPageId, widget);
    }
}

/**
 * Validates and persists widget resource changes to the API.
 */
async function saveWidget() {
    const id = elements.widgetId.value;
    const pageId = elements.widgetPageId.value;
    const type = elements.widgetType.value;
    const position = elements.widgetPosition.value ? parseInt(elements.widgetPosition.value) : 0;
    const configStr = elements.widgetConfig.value.trim();

    if (!type) {
        showToast('Please select a widget type', 'warning');
        return;
    }

    let config = {};
    if (configStr) {
        try {
            config = JSON.parse(configStr);
        } catch (e) {
            showToast('Invalid JSON in configuration', 'warning');
            return;
        }
    }

    const data = { type, position, config };

    try {
        if (id) {
            await apiCall(`/widgets/${id}`, 'PUT', data);
            showToast('Widget updated successfully', 'success');
        } else {
            await apiCall(`/pages/${pageId}/widgets`, 'POST', data);
            showToast('Widget created successfully', 'success');
        }

        closeWidgetModal();
        loadWidgets(pageId);
        loadDashboardData();
    } catch (error) {
        console.error('Failed to save widget:', error);
        showToast(error.error?.message || 'Failed to save widget', 'error');
    }
}

/**
 * Facilitates the deletion of a widget resource after user confirmation.
 * @param {string} widgetId - The identifier of the widget to delete.
 */
async function deleteWidget(widgetId) {
    if (!confirm('Are you sure you want to delete this widget?')) {
        return;
    }

    try {
        await apiCall(`/widgets/${widgetId}`, 'DELETE');
        showToast('Widget deleted successfully', 'success');
        loadWidgets(state.selectedPageId);
        loadDashboardData();
    } catch (error) {
        console.error('Failed to delete widget:', error);
        showToast(error.error?.message || 'Failed to delete widget', 'error');
    }
}

/**
 * Displays a transient notification message to the user.
 * @param {string} message - Content of the notification.
 * @param {string} [type='info'] - Severity level (success, error, warning, info).
 * @param {string} [title=null] - Optional notification title.
 */
function showToast(message, type = 'info', title = null) {
    const toast = document.createElement('div');
    toast.className = `toast toast-${type}`;

    const icons = {
        success: `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/>
            <polyline points="22 4 12 14.01 9 11.01"/>
        </svg>`,
        error: `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"/>
            <line x1="15" y1="9" x2="9" y2="15"/>
            <line x1="9" y1="9" x2="15" y2="15"/>
        </svg>`,
        warning: `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"/>
            <line x1="12" y1="9" x2="12" y2="13"/>
            <line x1="12" y1="17" x2="12.01" y2="17"/>
        </svg>`,
        info: `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"/>
            <line x1="12" y1="16" x2="12" y2="12"/>
            <line x1="12" y1="8" x2="12.01" y2="8"/>
        </svg>`
    };

    const titles = {
        success: 'Success',
        error: 'Error',
        warning: 'Warning',
        info: 'Info'
    };

    toast.innerHTML = `
        <span class="toast-icon">${icons[type] || icons.info}</span>
        <div class="toast-content">
            <div class="toast-title">${title || titles[type] || 'Notification'}</div>
            <div class="toast-message">${message}</div>
        </div>
        <button class="toast-close" onclick="this.parentElement.remove()">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="18" y1="6" x2="6" y2="18"/>
                <line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
        </button>
    `;

    elements.toastContainer.appendChild(toast);

    setTimeout(() => {
        if (toast.parentElement) {
            toast.remove();
        }
    }, 5000);
}

/**
 * Sanitizes input text by escaping HTML special characters.
 * @param {string} text - The raw text to sanitize.
 * @returns {string} Sanitized HTML string.
 */
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

/**
 * Formats an ISO date string into a localized short format.
 * @param {string} dateString - The ISO date string to format.
 * @returns {string} Localized date string.
 */
function formatDate(dateString) {
    if (!dateString) return 'N/A';
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', {
        month: 'short', day: 'numeric', year: 'numeric'
    });
}

window.openPageModal = openPageModal;
window.closePageModal = closePageModal;
window.savePage = savePage;
window.editPage = editPage;
window.deletePage = deletePage;
window.viewPage = viewPage;
window.viewPageWidgets = viewPageWidgets;
window.openWidgetModal = openWidgetModal;
window.closeWidgetModal = closeWidgetModal;
window.saveWidget = saveWidget;
window.editWidget = editWidget;
window.deleteWidget = deleteWidget;
window.goToPage = goToPage;

