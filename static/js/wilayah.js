(function () {
    const API_KOTA = "/api/v1/kota";
    const API_KECAMATAN = "/api/v1/kecamatan";

    const tbodyKota = document.querySelector("#tableKota tbody");
    const tbodyKecamatan = document.querySelector("#tableKecamatan tbody");
    const selectKota = document.querySelector("#kecamatanKotaId");


    // Fetch Kotas
    async function loadKotas() {
        try {
            const res = await fetch(API_KOTA, { credentials: "same-origin" });
            const json = await res.json();
            
            tbodyKota.innerHTML = "";
            selectKota.innerHTML = '<option value="">-- Pilih Kota --</option>';

            if (json.data && json.data.length > 0) {
                json.data.forEach((k, index) => {
                    // Update table
                    tbodyKota.innerHTML += `
                        <tr>
                            <td>${index + 1}</td>
                            <td>${k.name}</td>
                            <td>
                                <button class="btn btn-sm btn-info waves-effect waves-light btn-edit-kota" data-id="${k.id}" data-name="${k.name}" title="Edit">
                                    <i class="bx bx-pencil"></i>
                                </button>
                                <button class="btn btn-sm btn-danger waves-effect waves-light btn-delete-kota" data-id="${k.id}" title="Hapus">
                                    <i class="bx bx-trash"></i>
                                </button>
                            </td>
                        </tr>
                    `;
                    // Update dropdown
                    selectKota.innerHTML += `<option value="${k.id}">${k.name}</option>`;
                });
            } else {
                tbodyKota.innerHTML = `<tr><td colspan="3" class="text-center">Belum ada data Kota</td></tr>`;
            }
        } catch (error) {
            console.error("Error loading kotas:", error);
        }
    }

    // Fetch Kecamatans
    async function loadKecamatans() {
        try {
            const res = await fetch(API_KECAMATAN, { credentials: "same-origin" });
            const json = await res.json();
            
            tbodyKecamatan.innerHTML = "";

            if (json.data && json.data.length > 0) {
                json.data.forEach((k, index) => {
                    tbodyKecamatan.innerHTML += `
                        <tr>
                            <td>${index + 1}</td>
                            <td>${k.kota.name}</td>
                            <td>${k.name}</td>
                            <td>
                                <button class="btn btn-sm btn-info waves-effect waves-light btn-edit-kecamatan" data-id="${k.id}" data-kota-id="${k.kota_id}" data-name="${k.name}" title="Edit">
                                    <i class="bx bx-pencil"></i>
                                </button>
                                <button class="btn btn-sm btn-danger waves-effect waves-light btn-delete-kecamatan" data-id="${k.id}" title="Hapus">
                                    <i class="bx bx-trash"></i>
                                </button>
                            </td>
                        </tr>
                    `;
                });
            } else {
                tbodyKecamatan.innerHTML = `<tr><td colspan="4" class="text-center">Belum ada data Kecamatan</td></tr>`;
            }
        } catch (error) {
            console.error("Error loading kecamatans:", error);
        }
    }

    // Event Delegation for all Buttons
    document.addEventListener("click", async function(e) {
        // Create Kota
        const btnCreateKota = e.target.closest('#btnCreateKota');
        if (btnCreateKota) {
            e.preventDefault();
            const name = document.getElementById("kotaName").value;
            if (!name) {
                Swal.fire('Peringatan', 'Nama kota wajib diisi!', 'warning');
                return;
            }

            try {
                Swal.fire({ title: 'Mohon Tunggu', html: 'Sedang memproses...', allowOutsideClick: false, didOpen: () => { Swal.showLoading() } });
                const res = await fetch(API_KOTA, {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    credentials: "same-origin",
                    body: JSON.stringify({ name: name })
                });
                if (res.ok) {
                    document.getElementById("kotaName").value = "";
                    await loadKotas();
                    await loadKecamatans();
                    Swal.fire('Berhasil', 'Berhasil menambahkan kota!', 'success');
                } else {
                    const errData = await res.json();
                    Swal.fire('Gagal', "Gagal menambahkan kota: " + (errData.message || errData.error || "Unknown Error"), 'error');
                }
            } catch (error) {
                console.error(error);
                Swal.fire('Error', 'Terjadi kesalahan jaringan.', 'error');
            }
        }

        // Create Kecamatan
        const btnCreateKecamatan = e.target.closest('#btnCreateKecamatan');
        if (btnCreateKecamatan) {
            e.preventDefault();
            const kotaId = document.getElementById("kecamatanKotaId").value;
            const name = document.getElementById("kecamatanName").value;
            if (!kotaId || !name) {
                Swal.fire('Peringatan', 'Pilih kota dan isi nama kecamatan!', 'warning');
                return;
            }

            try {
                Swal.fire({ title: 'Mohon Tunggu', html: 'Sedang memproses...', allowOutsideClick: false, didOpen: () => { Swal.showLoading() } });
                const res = await fetch(API_KECAMATAN, {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    credentials: "same-origin",
                    body: JSON.stringify({ kota_id: parseInt(kotaId), name: name })
                });
                if (res.ok) {
                    document.getElementById("kecamatanName").value = "";
                    await loadKecamatans();
                    Swal.fire('Berhasil', 'Berhasil menambahkan kecamatan!', 'success');
                } else {
                    const errData = await res.json();
                    Swal.fire('Gagal', "Gagal menambahkan kecamatan: " + (errData.message || errData.error || "Unknown Error"), 'error');
                }
            } catch (error) {
                console.error(error);
                Swal.fire('Error', 'Terjadi kesalahan jaringan.', 'error');
            }
        }
        // Delete Kota
        const btnDeleteKota = e.target.closest('.btn-delete-kota');
        if (btnDeleteKota) {
            const id = btnDeleteKota.getAttribute('data-id');
            const confirm = await Swal.fire({
                title: 'Yakin ingin menghapus?',
                text: "Semua Kecamatan di dalamnya juga akan terhapus!",
                icon: 'warning',
                showCancelButton: true,
                confirmButtonColor: '#d33',
                cancelButtonColor: '#3085d6',
                confirmButtonText: 'Ya, hapus!',
                cancelButtonText: 'Batal'
            });
            if (!confirm.isConfirmed) return;

            try {
                Swal.fire({ title: 'Mohon Tunggu', html: 'Sedang memproses...', allowOutsideClick: false, didOpen: () => { Swal.showLoading() } });
                const res = await fetch(`${API_KOTA}/${id}`, { 
                    method: "DELETE",
                    credentials: "same-origin"
                });
                if (res.ok) {
                    await loadKotas();
                    await loadKecamatans();
                    Swal.fire('Terhapus!', 'Data kota berhasil dihapus.', 'success');
                } else {
                    Swal.fire('Gagal', 'Gagal menghapus kota', 'error');
                }
            } catch (error) {
                console.error(error);
                Swal.fire('Error', 'Terjadi kesalahan jaringan.', 'error');
            }
        }

        // Delete Kecamatan
        const btnDeleteKec = e.target.closest('.btn-delete-kecamatan');
        if (btnDeleteKec) {
            const id = btnDeleteKec.getAttribute('data-id');
            const confirm = await Swal.fire({
                title: 'Yakin ingin menghapus?',
                text: "Kecamatan ini akan dihapus secara permanen!",
                icon: 'warning',
                showCancelButton: true,
                confirmButtonColor: '#d33',
                cancelButtonColor: '#3085d6',
                confirmButtonText: 'Ya, hapus!',
                cancelButtonText: 'Batal'
            });
            if (!confirm.isConfirmed) return;

            try {
                Swal.fire({ title: 'Mohon Tunggu', html: 'Sedang memproses...', allowOutsideClick: false, didOpen: () => { Swal.showLoading() } });
                const res = await fetch(`${API_KECAMATAN}/${id}`, { 
                    method: "DELETE",
                    credentials: "same-origin"
                });
                if (res.ok) {
                    await loadKecamatans();
                    Swal.fire('Terhapus!', 'Data kecamatan berhasil dihapus.', 'success');
                } else {
                    Swal.fire('Gagal', 'Gagal menghapus kecamatan', 'error');
                }
            } catch (error) {
                console.error(error);
                Swal.fire('Error', 'Terjadi kesalahan jaringan.', 'error');
            }
        }

        // Edit Kota
        const btnEditKota = e.target.closest('.btn-edit-kota');
        if (btnEditKota) {
            const id = btnEditKota.getAttribute('data-id');
            const currentName = btnEditKota.getAttribute('data-name');
            
            const { value: newName } = await Swal.fire({
                title: 'Edit Nama Kota',
                input: 'text',
                inputValue: currentName,
                showCancelButton: true,
                inputValidator: (value) => {
                    if (!value) return 'Nama kota wajib diisi!'
                }
            });

            if (newName && newName !== currentName) {
                try {
                    const res = await fetch(`${API_KOTA}/${id}`, {
                        method: "PUT",
                        headers: { "Content-Type": "application/json" },
                        credentials: "same-origin",
                        body: JSON.stringify({ name: newName })
                    });
                    if (res.ok) {
                        await loadKotas();
                        await loadKecamatans();
                        Swal.fire('Tersimpan!', 'Nama kota berhasil diubah.', 'success');
                    } else {
                        Swal.fire('Gagal', 'Gagal mengubah kota', 'error');
                    }
                } catch (error) {
                    console.error(error);
                    Swal.fire('Error', 'Terjadi kesalahan jaringan.', 'error');
                }
            }
        }

        // Edit Kecamatan
        const btnEditKecamatan = e.target.closest('.btn-edit-kecamatan');
        if (btnEditKecamatan) {
            const id = btnEditKecamatan.getAttribute('data-id');
            const currentName = btnEditKecamatan.getAttribute('data-name');
            const currentKotaId = btnEditKecamatan.getAttribute('data-kota-id');
            
            // Get kota options for the select element
            const kotaOptions = Array.from(selectKota.options)
                .filter(opt => opt.value !== "")
                .map(opt => `<option value="${opt.value}" ${opt.value === currentKotaId ? 'selected' : ''}>${opt.text}</option>`)
                .join('');

            const { value: formValues } = await Swal.fire({
                title: 'Edit Kecamatan',
                html:
                    `<select id="swal-kotaId" class="swal2-select" style="display: flex; margin: 10px auto; width: 80%;">${kotaOptions}</select>` +
                    `<input id="swal-kecName" class="swal2-input" value="${currentName}">`,
                focusConfirm: false,
                showCancelButton: true,
                preConfirm: () => {
                    const kotaId = document.getElementById('swal-kotaId').value;
                    const name = document.getElementById('swal-kecName').value;
                    if (!kotaId || !name) {
                        Swal.showValidationMessage('Semua kolom wajib diisi!');
                    }
                    return { kotaId, name };
                }
            });

            if (formValues) {
                try {
                    Swal.fire({ title: 'Mohon Tunggu', html: 'Sedang memproses...', allowOutsideClick: false, didOpen: () => { Swal.showLoading() } });
                    const res = await fetch(`${API_KECAMATAN}/${id}`, {
                        method: "PUT",
                        headers: { "Content-Type": "application/json" },
                        credentials: "same-origin",
                        body: JSON.stringify({ kota_id: parseInt(formValues.kotaId), name: formValues.name })
                    });
                    if (res.ok) {
                        await loadKecamatans();
                        Swal.fire('Tersimpan!', 'Kecamatan berhasil diubah.', 'success');
                    } else {
                        Swal.fire('Gagal', 'Gagal mengubah kecamatan', 'error');
                    }
                } catch (error) {
                    console.error(error);
                    Swal.fire('Error', 'Terjadi kesalahan jaringan.', 'error');
                }
            }
        }
    });

    // Initial load
    loadKotas().then(() => loadKecamatans());
})();
