let currentPage = 1;
const pageSize = 10;

async function loadItems() {
    const categoryFilter = document.getElementById('categoryFilter').value;
    const url = `/api/items?page=${currentPage}&page_size=${pageSize}${categoryFilter ? '&category_type=' + categoryFilter : ''}`;
    
    try {
        const response = await fetch(url);
        const result = await response.json();
        
        if (result.success) {
            displayItems(result.data.items);
            updatePagination(result.data.total, result.data.page, result.data.page_size);
        } else {
            alert('エラー: ' + result.error);
        }
    } catch (error) {
        alert('データの取得に失敗しました: ' + error.message);
    }
}

function displayItems(items) {
    const tbody = document.getElementById('itemsTableBody');
    tbody.innerHTML = '';
    
    items.forEach(item => {
        const row = document.createElement('tr');
        
        let detailAttributes = '';
        if (item.type_a) {
            detailAttributes = `容量: ${item.type_a.capacity || '-'}<br>材質: ${item.type_a.material || '-'}`;
        } else if (item.type_b) {
            detailAttributes = `内径: ${item.type_b.inner_diameter || '-'}<br>外径: ${item.type_b.outer_diameter || '-'}`;
        }
        
        row.innerHTML = `
            <td>${item.item_id}</td>
            <td>${item.item_name}</td>
            <td>${item.category_type}</td>
            <td>${item.item_code}</td>
            <td class="detail-attributes">${detailAttributes}</td>
            <td>
                <div class="action-buttons">
                    <button onclick="showEditForm('${item.item_id}')" class="btn btn-warning">編集</button>
                    <button onclick="deleteItem('${item.item_id}')" class="btn btn-danger">削除</button>
                </div>
            </td>
        `;
        
        tbody.appendChild(row);
    });
}

function updatePagination(total, page, pageSize) {
    const totalPages = Math.ceil(total / pageSize);
    const pagination = document.getElementById('pagination');
    
    pagination.innerHTML = `
        <button onclick="changePage(${page - 1})" ${page === 1 ? 'disabled' : ''}>前へ</button>
        <span class="current-page">ページ ${page} / ${totalPages}</span>
        <button onclick="changePage(${page + 1})" ${page === totalPages ? 'disabled' : ''}>次へ</button>
    `;
}

function changePage(page) {
    currentPage = page;
    loadItems();
}

function showCreateForm() {
    document.getElementById('createForm').style.display = 'block';
    document.getElementById('editForm').style.display = 'none';
}

function hideCreateForm() {
    document.getElementById('createForm').style.display = 'none';
    document.getElementById('createForm').reset();
    document.getElementById('categoryAFields').style.display = 'none';
    document.getElementById('categoryBFields').style.display = 'none';
}

function toggleCategoryFields() {
    const categoryType = document.getElementById('categoryType').value;
    document.getElementById('categoryAFields').style.display = categoryType === 'A' ? 'block' : 'none';
    document.getElementById('categoryBFields').style.display = categoryType === 'B' ? 'block' : 'none';
}

async function createItem(event) {
    event.preventDefault();
    
    const categoryType = document.getElementById('categoryType').value;
    const data = {
        item_id: document.getElementById('itemId').value,
        item_name: document.getElementById('itemName').value,
        category_type: categoryType,
        item_code: document.getElementById('itemCode').value
    };
    
    if (categoryType === 'A') {
        const capacity = document.getElementById('capacity').value;
        const material = document.getElementById('material').value;
        if (capacity) data.capacity = parseFloat(capacity);
        if (material) data.material = material;
    } else if (categoryType === 'B') {
        const innerDiameter = document.getElementById('innerDiameter').value;
        const outerDiameter = document.getElementById('outerDiameter').value;
        if (innerDiameter) data.inner_diameter = parseFloat(innerDiameter);
        if (outerDiameter) data.outer_diameter = parseFloat(outerDiameter);
    }
    
    try {
        const response = await fetch('/api/items', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        });
        
        const result = await response.json();
        
        if (result.success) {
            alert('品目を追加しました');
            hideCreateForm();
            loadItems();
        } else {
            alert('エラー: ' + result.error);
        }
    } catch (error) {
        alert('品目の追加に失敗しました: ' + error.message);
    }
}

async function showEditForm(itemId) {
    try {
        const response = await fetch(`/api/items/${itemId}`);
        const result = await response.json();
        
        if (result.success) {
            const item = result.data;
            
            document.getElementById('editItemId').value = item.item_id;
            document.getElementById('editCategoryType').value = item.category_type;
            document.getElementById('editItemName').value = item.item_name;
            document.getElementById('editItemCode').value = item.item_code;
            
            if (item.category_type === 'A' && item.type_a) {
                document.getElementById('editCapacity').value = item.type_a.capacity || '';
                document.getElementById('editMaterial').value = item.type_a.material || '';
                document.getElementById('editCategoryAFields').style.display = 'block';
                document.getElementById('editCategoryBFields').style.display = 'none';
            } else if (item.category_type === 'B' && item.type_b) {
                document.getElementById('editInnerDiameter').value = item.type_b.inner_diameter || '';
                document.getElementById('editOuterDiameter').value = item.type_b.outer_diameter || '';
                document.getElementById('editCategoryAFields').style.display = 'none';
                document.getElementById('editCategoryBFields').style.display = 'block';
            }
            
            document.getElementById('createForm').style.display = 'none';
            document.getElementById('editForm').style.display = 'block';
        } else {
            alert('エラー: ' + result.error);
        }
    } catch (error) {
        alert('品目情報の取得に失敗しました: ' + error.message);
    }
}

function hideEditForm() {
    document.getElementById('editForm').style.display = 'none';
}

async function updateItem(event) {
    event.preventDefault();
    
    const itemId = document.getElementById('editItemId').value;
    const categoryType = document.getElementById('editCategoryType').value;
    const data = {};
    
    const itemName = document.getElementById('editItemName').value;
    const itemCode = document.getElementById('editItemCode').value;
    
    if (itemName) data.item_name = itemName;
    if (itemCode) data.item_code = itemCode;
    
    if (categoryType === 'A') {
        const capacity = document.getElementById('editCapacity').value;
        const material = document.getElementById('editMaterial').value;
        if (capacity) data.capacity = parseFloat(capacity);
        if (material) data.material = material;
    } else if (categoryType === 'B') {
        const innerDiameter = document.getElementById('editInnerDiameter').value;
        const outerDiameter = document.getElementById('editOuterDiameter').value;
        if (innerDiameter) data.inner_diameter = parseFloat(innerDiameter);
        if (outerDiameter) data.outer_diameter = parseFloat(outerDiameter);
    }
    
    try {
        const response = await fetch(`/api/items/${itemId}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        });
        
        const result = await response.json();
        
        if (result.success) {
            alert('品目を更新しました');
            hideEditForm();
            loadItems();
        } else {
            alert('エラー: ' + result.error);
        }
    } catch (error) {
        alert('品目の更新に失敗しました: ' + error.message);
    }
}

async function deleteItem(itemId) {
    if (!confirm(`品目ID ${itemId} を削除してもよろしいですか？`)) {
        return;
    }
    
    try {
        const response = await fetch(`/api/items/${itemId}`, {
            method: 'DELETE'
        });
        
        const result = await response.json();
        
        if (result.success) {
            alert('品目を削除しました');
            loadItems();
        } else {
            alert('エラー: ' + result.error);
        }
    } catch (error) {
        alert('品目の削除に失敗しました: ' + error.message);
    }
}

// ページ読み込み時に品目一覧を表示
window.addEventListener('DOMContentLoaded', () => {
    loadItems();
});