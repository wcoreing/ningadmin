<template>
  <!-- 图片预览组件 -->
  <el-image-viewer
    v-if="previewVisible"
    :url-list="previewUrlList"
    :initial-index="previewIndex"
    :hide-on-click-modal="false"
    @close="previewVisible = false"
  />

  <!-- 上传文件弹窗 -->
  <el-dialog
    v-model="uploadVisible"
    title="上传文件"
    width="500px"
    :append-to-body="true"
    :destroy-on-close="false"
    :close-on-click-modal="false"
    :show-close="true"
    @close="handleDialogClose"
  >
    <el-upload
      class="file-uploader"
      :http-request="customUpload"
      :show-file-list="true"
      :on-success="handleUploadSuccess"
      :on-error="handleUploadError"
      :before-upload="beforeUpload"
      multiple
      drag
    >
      <el-icon class="el-icon--upload"><upload-filled /></el-icon>
      <div class="el-upload__text">
        将文件拖到此处，或<em>点击上传</em>
      </div>
      <template #tip>
        <div class="el-upload__tip">
          支持常见图片格式，单个文件不超过 10MB
        </div>
      </template>
    </el-upload>
  </el-dialog>

  <template v-if="standalone">
    <div class="file-manager" @paste="handlePaste">
      <!-- 左侧分类树 -->
      <div class="category-tree">
        <div class="tree-header">
          <el-input
            v-model="categoryKeyword"
            placeholder="搜索分类"
            size="small"
            :prefix-icon="Search"
            clearable
          />
        </div>
        <el-scrollbar>
          <el-tree
            ref="treeRef"
            :data="categoryTree"
            :props="{ label: 'name' }"
            :filter-node-method="filterNode"
            @node-click="handleCategoryClick"
            default-expand-all
            highlight-current
          >
            <template #default="{ node, data }">
              <div class="category-node">
                <el-icon><component :is="data.icon || 'Folder'" /></el-icon>
                <span>{{ node.label }}</span>
                <el-tag size="small" type="info" class="count">{{ data.count || 0 }}</el-tag>
              </div>
            </template>
          </el-tree>
        </el-scrollbar>
      </div>

      <!-- 右侧内容区 -->
      <div class="content-area">
        <!-- 添加粘贴提示 -->
        <div class="paste-tip">
          <el-alert
            title="支持直接粘贴图片文件(Ctrl+V)"
            type="info"
            :closable="false"
            show-icon
          />
        </div>

        <!-- 顶部工具栏 -->
        <div class="toolbar">
          <div class="left">
            <el-button-group class="main-actions">
              <el-button 
                type="primary" 
                @click="handleUpload"
                ref="uploadBtn"
              >
                <el-icon><Upload /></el-icon>上传文件
              </el-button>
              <el-button 
                @click="handleBatchMove" 
                :disabled="!selectedFiles.length"
              >
                <el-icon><FolderAdd /></el-icon>移动到分类
              </el-button>
              <el-button @click="handleRemove" :disabled="!selectedFiles.length">
                <el-icon><Delete /></el-icon>删除
              </el-button>
              <el-button @click="handleDownload" :disabled="!selectedFiles.length">
                <el-icon><Download /></el-icon>下载
              </el-button>
            </el-button-group>
          </div>

          <div class="right">
            <el-input
              v-model="searchKeyword"
              placeholder="搜索文件名"
              class="search-input"
              :prefix-icon="Search"
              clearable
              @input="handleSearch"
            />
            <el-select v-model="currentSort" class="sort-select" size="default">
              <el-option
                v-for="(label, key) in sortLabels"
                :key="key"
                :label="label"
                :value="key"
              />
            </el-select>
            <el-radio-group v-model="viewMode" size="default">
              <el-tooltip content="网格视图" placement="top">
                <el-radio-button label="grid">
                  <el-icon><Grid /></el-icon>
                </el-radio-button>
              </el-tooltip>
              <el-tooltip content="列表视图" placement="top">
                <el-radio-button label="list">
                  <el-icon><List /></el-icon>
                </el-radio-button>
              </el-tooltip>
            </el-radio-group>
          </div>
        </div>

        <!-- 快捷操作栏 -->
        <div class="quick-actions">
          <div class="selection-actions">
            <el-button-group>
              <el-button size="small" @click="handleSelectAll">全选</el-button>
              <el-button size="small" @click="handleInvertSelection">反选</el-button>
              <el-button size="small" @click="handleClearSelection">清空</el-button>
            </el-button-group>
          </div>
          <div class="selected-count" v-if="selectedFiles.length">
            已选择 {{ selectedFiles.length }} 个文件
          </div>
        </div>

        <!-- 文件列表区域 -->
        <el-scrollbar class="file-list-container">
          <!-- 网格视图 -->
          <div v-if="viewMode === 'grid'" class="file-grid">
            <div
              v-for="item in fileList"
              :key="item.id"
              class="file-item"
              :class="{ selected: selectedFiles.includes(item.id) }"
            >
              <!-- 选择遮罩 -->
              <div class="select-mask" @click="handleSelectFile(item)">
                <el-icon v-if="selectedFiles.includes(item.id)"><Select /></el-icon>
              </div>

              <!-- 文件预览 -->
              <div class="file-preview">
                <el-image
                  v-if="isImage(item.type)"
                  :src="item.url"
                  fit="contain"
                  :preview-src-list="[item.url]"
                  :preview-teleported="true"
                  :hide-on-click-modal="false"
                  preview-class="image-preview"
                />
                <div v-else class="file-icon">
                  <el-icon><Document /></el-icon>
                </div>
              </div>

              <!-- 文件信息 -->
              <div class="file-info">
                <div class="file-info-content">
                  <div class="file-name" :title="item.originalName">{{ item.originalName }}</div>
                  <div class="system-name" :title="`ID: ${item.id} | ${item.name}`">ID: {{ item.id }} | {{ item.name }}</div>
                </div>
                <div class="file-meta">
                  <span class="size">{{ formatSize(item.size) }}</span>
                  <span class="date">{{ formatDate(item.createdAt) }}</span>
                </div>
              </div>

              <!-- 操作按钮 -->
              <div class="file-actions">
                <template v-if="isImage(item.type)">
                  <el-tooltip content="预览" placement="top">
                    <el-button circle @click.stop="handleImagePreview(item)">
                      <el-icon><ZoomIn /></el-icon>
                    </el-button>
                  </el-tooltip>
                </template>
                <el-tooltip content="下载" placement="top">
                  <el-button circle @click.stop="handleDownloadFile(item)">
                    <el-icon><Download /></el-icon>
                  </el-button>
                </el-tooltip>
                <el-tooltip content="删除" placement="top">
                  <el-button circle @click.stop="handleDeleteFile(item)">
                    <el-icon><Delete /></el-icon>
                  </el-button>
                </el-tooltip>
              </div>
            </div>
          </div>

          <!-- 列表视图 -->
          <el-table
            v-else
            :data="fileList"
            style="width: 100%"
            @selection-change="handleTableSelection"
          >
            <el-table-column type="selection" width="55" />
            <el-table-column label="文件名" min-width="200">
              <template #default="{ row }">
                <div class="file-name-cell">
                  <el-icon v-if="isImage(row.type)"><Picture /></el-icon>
                  <el-icon v-else><Document /></el-icon>
                  <div class="file-name-content">
                    <div class="original-name">{{ row.originalName }}</div>
                    <div class="system-name">ID: {{ row.id }} | {{ row.name }}</div>
                  </div>
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="size" label="大小" width="120">
              <template #default="{ row }">
                {{ formatSize(row.size) }}
              </template>
            </el-table-column>
            <el-table-column prop="createdAt" label="上传时间" width="180">
              <template #default="{ row }">
                {{ formatDate(row.createdAt) }}
              </template>
            </el-table-column>
            <el-table-column label="操作" width="180">
              <template #default="{ row }">
                <el-button-group>
                  <el-button link type="primary" @click="handlePreview(row)">
                    <el-icon><View /></el-icon>
                  </el-button>
                  <el-button link type="primary" @click="handleDownloadFile(row)">
                    <el-icon><Download /></el-icon>
                  </el-button>
                  <el-button link type="danger" @click="handleDeleteFile(row)">
                    <el-icon><Delete /></el-icon>
                  </el-button>
                </el-button-group>
              </template>
            </el-table-column>
          </el-table>
        </el-scrollbar>

        <!-- 分页 -->
        <div class="pagination">
          <el-pagination
            v-model:current-page="page"
            v-model:page-size="pageSize"
            :total="total"
            :page-sizes="props.standalone ? [55, 110, 165, 220] : [55, 110, 165, 220]"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
          />
        </div>
      </div>

      <!-- 移动到分类的对话框 -->
      <el-dialog
        v-model="moveDialogVisible"
        title="移动到分类"
        width="400px"
      >
        <el-tree
          ref="moveTreeRef"
          :data="categoryTree"
          :props="{ label: 'name' }"
          node-key="id"
          default-expand-all
          highlight-current
          @node-click="handleMoveToCategory"
        >
          <template #default="{ node, data }">
            <div class="category-node">
              <el-icon><component :is="data.icon || 'Folder'" /></el-icon>
              <span>{{ node.label }}</span>
            </div>
          </template>
        </el-tree>
        <template #footer>
          <div class="dialog-footer">
            <el-button @click="moveDialogVisible = false">取消</el-button>
            <el-button type="primary" @click="moveDialogVisible = false">
              确定
            </el-button>
          </div>
        </template>
      </el-dialog>
    </div>
  </template>
  <template v-else>
    <el-dialog
      v-model="dialogVisible"
      title="文件管理"
      width="70%"
      :fullscreen="false"
      :close-on-click-modal="false"
      :close-on-press-escape="false"
      class="file-manager-dialog"
    >
      <div class="file-manager" @paste="handlePaste">
        <!-- 左侧分类树 -->
        <div class="category-tree">
          <div class="tree-header">
            <el-input
              v-model="categoryKeyword"
              placeholder="搜索分类"
              size="small"
              :prefix-icon="Search"
              clearable
            />
          </div>
          <el-scrollbar>
            <el-tree
              ref="treeRef"
              :data="categoryTree"
              :props="{ label: 'name' }"
              :filter-node-method="filterNode"
              @node-click="handleCategoryClick"
              default-expand-all
              highlight-current
            >
              <template #default="{ node, data }">
                <div class="category-node">
                  <el-icon><component :is="data.icon || 'Folder'" /></el-icon>
                  <span>{{ node.label }}</span>
                  <el-tag size="small" type="info" class="count">{{ data.count || 0 }}</el-tag>
                </div>
              </template>
            </el-tree>
          </el-scrollbar>
        </div>

        <!-- 右侧内容区 -->
        <div class="content-area">
          <!-- 添加粘贴提示 -->
          <div class="paste-tip">
            <el-alert
              title="支持直接粘贴图片文件(Ctrl+V)"
              type="info"
              :closable="false"
              show-icon
            />
          </div>

          <!-- 顶部工具栏 -->
          <div class="toolbar">
            <div class="left">
              <el-button-group class="main-actions">
                <el-button 
                  type="primary" 
                  @click="handleUpload"
                  ref="uploadBtn"
                >
                  <el-icon><Upload /></el-icon>上传文件
                </el-button>
                <el-button 
                  @click="handleBatchMove" 
                  :disabled="!selectedFiles.length"
                >
                  <el-icon><FolderAdd /></el-icon>移动到分类
                </el-button>
                <el-button @click="handleRemove" :disabled="!selectedFiles.length">
                  <el-icon><Delete /></el-icon>删除
                </el-button>
              </el-button-group>
            </div>

            <div class="right">
              <el-input
                v-model="searchKeyword"
                placeholder="搜索文件名"
                class="search-input"
                :prefix-icon="Search"
                clearable
                @input="handleSearch"
              />
              <el-select v-model="currentSort" class="sort-select" size="default">
                <el-option
                  v-for="(label, key) in sortLabels"
                  :key="key"
                  :label="label"
                  :value="key"
                />
              </el-select>
              <el-radio-group v-model="viewMode" size="default">
                <el-tooltip content="网格视图" placement="top">
                  <el-radio-button label="grid">
                    <el-icon><Grid /></el-icon>
                  </el-radio-button>
                </el-tooltip>
                <el-tooltip content="列表视图" placement="top">
                  <el-radio-button label="list">
                    <el-icon><List /></el-icon>
                  </el-radio-button>
                </el-tooltip>
              </el-radio-group>
            </div>
          </div>

          <!-- 快捷操作栏 -->
          <div class="quick-actions">
            <div class="selection-actions">
              <el-button-group>
                <el-button size="small" @click="handleSelectAll">全选</el-button>
                <el-button size="small" @click="handleInvertSelection">反选</el-button>
                <el-button size="small" @click="handleClearSelection">清空</el-button>
              </el-button-group>
            </div>
            <div class="selected-count" v-if="selectedFiles.length">
              已选择 {{ selectedFiles.length }} 个文件
            </div>
          </div>

          <!-- 文件列表区域 -->
          <el-scrollbar class="file-list-container">
            <!-- 网格视图 -->
            <div v-if="viewMode === 'grid'" class="file-grid">
              <div
                v-for="item in fileList"
                :key="item.id"
                class="file-item"
                :class="{ selected: selectedFiles.includes(item.id) }"
              >
                <!-- 选择遮罩 -->
                <div class="select-mask" @click="handleSelectFile(item)">
                  <el-icon v-if="selectedFiles.includes(item.id)"><Select /></el-icon>
                </div>

                <!-- 文件预览 -->
                <div class="file-preview">
                  <el-image
                    v-if="isImage(item.type)"
                    :src="item.url"
                    fit="contain"
                    :preview-src-list="[item.url]"
                    :preview-teleported="true"
                    :hide-on-click-modal="false"
                    preview-class="image-preview"
                  />
                  <div v-else class="file-icon">
                    <el-icon><Document /></el-icon>
                  </div>
                </div>

                <!-- 文件信息 -->
                <div class="file-info">
                  <div class="file-info-content">
                    <div class="file-name" :title="item.originalName">{{ item.originalName }}</div>
                    <div class="system-name" :title="`ID: ${item.id} | ${item.name}`">ID: {{ item.id }} | {{ item.name }}</div>
                  </div>
                  <div class="file-meta">
                    <span class="size">{{ formatSize(item.size) }}</span>
                    <span class="date">{{ formatDate(item.createdAt) }}</span>
                  </div>
                </div>

                <!-- 操作按钮 -->
                <div class="file-actions">
                  <template v-if="isImage(item.type)">
                    <el-tooltip content="预览" placement="top">
                      <el-button circle @click.stop="handleImagePreview(item)">
                        <el-icon><ZoomIn /></el-icon>
                      </el-button>
                    </el-tooltip>
                  </template>
                  <el-tooltip content="下载" placement="top">
                    <el-button circle @click.stop="handleDownloadFile(item)">
                      <el-icon><Download /></el-icon>
                    </el-button>
                  </el-tooltip>
                  <el-tooltip content="删除" placement="top">
                    <el-button circle @click.stop="handleDeleteFile(item)">
                      <el-icon><Delete /></el-icon>
                    </el-button>
                  </el-tooltip>
                </div>
              </div>
            </div>

            <!-- 列表视图 -->
            <el-table
              v-else
              :data="fileList"
              style="width: 100%"
              @selection-change="handleTableSelection"
            >
              <el-table-column type="selection" width="55" />
              <el-table-column label="文件名" min-width="200">
                <template #default="{ row }">
                  <div class="file-name-cell">
                    <el-icon v-if="isImage(row.type)"><Picture /></el-icon>
                    <el-icon v-else><Document /></el-icon>
                    <div class="file-name-content">
                      <div class="original-name">{{ row.originalName }}</div>
                      <div class="system-name">ID: {{ row.id }} | {{ row.name }}</div>
                    </div>
                  </div>
                </template>
              </el-table-column>
              <el-table-column prop="size" label="大小" width="120">
                <template #default="{ row }">
                  {{ formatSize(row.size) }}
                </template>
              </el-table-column>
              <el-table-column prop="createdAt" label="上传时间" width="180">
                <template #default="{ row }">
                  {{ formatDate(row.createdAt) }}
                </template>
              </el-table-column>
              <el-table-column label="操作" width="180">
                <template #default="{ row }">
                  <el-button-group>
                    <el-button link type="primary" @click="handlePreview(row)">
                      <el-icon><View /></el-icon>
                    </el-button>
                    <el-button link type="primary" @click="handleDownloadFile(row)">
                      <el-icon><Download /></el-icon>
                    </el-button>
                    <el-button link type="danger" @click="handleDeleteFile(row)">
                      <el-icon><Delete /></el-icon>
                    </el-button>
                  </el-button-group>
                </template>
              </el-table-column>
            </el-table>
          </el-scrollbar>

          <!-- 分页 -->
          <div class="pagination">
            <el-pagination
              v-model:current-page="page"
              v-model:page-size="pageSize"
              :total="total"
              :page-sizes="props.standalone ? [55, 110, 165, 220] : [55, 110, 165, 220]"
              layout="total, sizes, prev, pager, next, jumper"
              @size-change="handleSizeChange"
              @current-change="handleCurrentChange"
            />
          </div>
        </div>
      </div>

      <!-- 移动到分类的对话框 -->
      <el-dialog
        v-model="moveVisible"
        title="移动到"
        width="400px"
        append-to-body
      >
        <el-tree
          :data="categoryTree"
          :props="{ label: 'name' }"
          @node-click="handleMoveToCategory"
          default-expand-all
        />
      </el-dialog>

      <!-- 添加底部操作按钮 -->
      <template #footer>
        <div class="dialog-footer">
          <div class="selected-info" v-if="selectedFiles.length">
            已选择 {{ selectedFiles.length }} 个文件
          </div>
          <div class="actions">
            <el-button @click="handleClose">取 消</el-button>
            <el-button
              type="primary"
              :disabled="!canConfirm"
              @click="handleConfirmSelect"
            >
              确认选择
            </el-button>
          </div>
        </div>
      </template>
    </el-dialog>
  </template>
</template>

<style lang="scss" scoped>
.file-manager-dialog {
  :deep(.el-dialog__body) {
    padding: 0;
    height: calc(60vh - 120px);
  }

  :deep(.el-dialog__header) {
    margin: 0;
    padding: 16px 20px;
    border-bottom: 1px solid var(--el-border-color-lighter);
  }
}

.file-manager {
  display: flex;
  height: 100%;
  background: var(--el-bg-color-page);

  // 左侧分类树样式
  .category-tree {
    width: 220px;
    border-right: 1px solid var(--el-border-color-lighter);
    display: flex;
    flex-direction: column;

    .tree-header {
      padding: 16px;
      border-bottom: 1px solid var(--el-border-color-lighter);
    }

    .el-scrollbar {
      flex: 1;
      padding: 12px;
    }

    .category-node {
      display: flex;
      align-items: center;
      gap: 8px;
      padding: 4px 0;

      .count {
        margin-left: auto;
      }
    }
  }

  // 右侧内容区样式
  .content-area {
    flex: 1;
    display: flex;
    flex-direction: column;
    background: #fff;

    // 工具栏样式
    .toolbar {
      padding: 16px;
      border-bottom: 1px solid var(--el-border-color-lighter);
      display: flex;
      justify-content: space-between;
      align-items: center;
      gap: 16px;

      .right {
        display: flex;
        align-items: center;
        gap: 12px;

        .search-input {
          width: 240px;
        }

        .sort-select {
          width: 120px;
        }
      }
    }

    // 快捷操作栏样式
    .quick-actions {
      padding: 8px 16px;
      background: var(--el-fill-color-light);
      display: flex;
      justify-content: space-between;
      align-items: center;

      .selected-count {
        color: var(--el-text-color-secondary);
        font-size: 14px;
      }
    }

    // 文件列表容器
    .file-list-container {
      flex: 1;
      overflow: hidden;
      padding: 16px;

      .file-grid {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
        gap: 16px;
        padding: 4px;

        .file-item {
          position: relative;
          border: 1px solid var(--el-border-color-lighter);
          border-radius: 4px;
          overflow: hidden;
          transition: all 0.2s;

          &:hover {
            transform: translateY(-2px);
            box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);

            .file-actions {
              opacity: 1;
            }
          }

          &.selected {
            border-color: var(--el-color-primary);
            background: var(--el-color-primary-light-9);

            .select-mask {
              opacity: 1;
            }
          }

          .select-mask {
            position: absolute;
            inset: 0;
            display: flex;
            align-items: center;
            justify-content: center;
            background: rgba(var(--el-color-primary-rgb), 0.1);
            opacity: 0;
            transition: opacity 0.2s;
            z-index: 1;
            cursor: pointer;

            .el-icon {
              font-size: 24px;
              color: var(--el-color-primary);
              background: #fff;
              border-radius: 50%;
              padding: 4px;
            }
          }

          .file-preview {
            height: 120px;
            display: flex;
            align-items: center;
            justify-content: center;
            background: var(--el-fill-color-lighter);

            .el-image {
              width: 100%;
              height: 100%;
              object-fit: contain;
            }

            .file-icon {
              font-size: 32px;
              color: var(--el-text-color-secondary);
            }
          }

          .file-info {
            padding: 8px;

            .file-info-content {
              font-size: 13px;
              line-height: 1.4;
              margin-bottom: 4px;
              
              .file-name {
                display: -webkit-box;
                -webkit-line-clamp: 1;
                -webkit-box-orient: vertical;
                overflow: hidden;
                text-overflow: ellipsis;
                word-break: break-all;
              }

              .system-name {
                font-size: 12px;
                color: var(--el-text-color-secondary);
                margin-top: 2px;
                display: -webkit-box;
                -webkit-line-clamp: 1;
                -webkit-box-orient: vertical;
                overflow: hidden;
                text-overflow: ellipsis;
              }
            }

            .file-meta {
              display: flex;
              justify-content: space-between;
              font-size: 12px;
              color: var(--el-text-color-secondary);
            }
          }

          .file-actions {
            position: absolute;
            right: 8px;
            bottom: 8px;
            display: flex;
            gap: 4px;
            opacity: 0;
            transition: opacity 0.2s;
            z-index: 2;

            .el-button {
              padding: 6px;
              background: rgba(255, 255, 255, 0.9);
              border: none;
              box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);

              &:hover {
                background: var(--el-color-primary);
                color: #fff;
                transform: scale(1.1);
              }
            }
          }
        }
      }
    }

    // 分页样式
    .pagination {
      padding: 12px 16px;
      border-top: 1px solid var(--el-border-color-lighter);
      display: flex;
      justify-content: center;
    }
  }
}

.dialog-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 20px;
  
  .selected-info {
    color: var(--el-text-color-secondary);
  }
  
  .actions {
    display: flex;
    gap: 10px;
  }
}

// 独立页面模式的样式
.file-manager {
  height: calc(100vh - 120px);
  background: var(--el-bg-color-page);
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 4px;
}

// 列表视图的文件名单元格样式
.file-name-cell {
  display: flex;
  align-items: center;
  gap: 8px;
  
  .file-name-content {
    flex: 1;
    min-width: 0;

    .original-name {
      font-size: 14px;
    }

    .system-name {
      font-size: 12px;
      color: var(--el-text-color-secondary);
      margin-top: 2px;
    }
  }
}

:deep(.image-preview) {
  .el-image-viewer__wrapper {
    .el-image-viewer__close {
      color: #fff;
    }
  }
}

// 添加粘贴提示
.paste-tip {
  padding: 4px 8px;
  margin-bottom: 0;
  background: #f5f7fa;
  
  :deep(.el-alert) {
    margin: 0;
    padding: 4px 8px;
    background: transparent;
    border: none;
    
    .el-alert__content {
      padding: 0 4px;
    }
    
    .el-alert__icon {
      font-size: 14px;
    }
    
    .el-alert__title {
      font-size: 12px;
      line-height: 1.2;
      color: var(--el-text-color-secondary);
    }
  }
}

// 移动分类对话框样式
.el-dialog {
  :deep(.el-tree) {
    background: none;
    
    .el-tree-node__content {
      height: 32px;
    }
  }
  
  .category-node {
    display: flex;
    align-items: center;
    gap: 8px;
    
    .el-icon {
      font-size: 16px;
      color: var(--el-text-color-secondary);
    }
  }
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

// 确保上传对话框样式在最外层
.el-dialog {
  :deep(.el-dialog__body) {
    padding: 20px;
  }
}

.file-uploader {
  :deep(.el-upload) {
    width: 100%;
  }
  
  :deep(.el-upload-dragger) {
    width: 100%;
    height: 200px;
  }
}
</style>

<script setup>
import { ref, onMounted, watch, computed, onUnmounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ElImageViewer } from 'element-plus'
import { Folder, Upload, Delete, Search, Document, UploadFilled, Grid, List, ArrowDown, FolderAdd, View, Check, Select, ZoomIn } from '@element-plus/icons-vue'
import { getFileList, deleteFile, getFileCategories, uploadFile, moveFiles } from '@/admin/api/common/upload'
import { getToken } from '@/admin/utils/auth'
import dayjs from 'dayjs'

const props = defineProps({
  standalone: {
    type: Boolean,
    default: false
  },
  modelValue: {
    type: Boolean,
    default: false
  },
  multiple: {
    type: Boolean,
    default: true
  },
  pageSize: {
    type: Number,
    default: 0  // 0 表示未设置
  }
})

const emit = defineEmits(['update:modelValue', 'select'])

// 弹窗显示控制
const dialogVisible = ref(false)
const uploadVisible = ref(false)
const moveVisible = ref(false)

// 文件列表相关
const page = ref(1)
const pageSize = ref(props.pageSize || (props.standalone ? 55 : 55))  // 设置为55条
const total = ref(0)
const searchKeyword = ref('')
const selectedFiles = ref([])
const fileList = ref([])

// 分类树数据
const categoryTree = ref([])

// 新增功能相关的变量
const viewMode = ref('grid')
const categoryKeyword = ref('')
const treeRef = ref(null)

// 排序相关
const currentSort = ref('name')
const sortLabels = {
  name: '按名称排序',
  date: '按时间排序',
  size: '按大小排序'
}

// 预览相关的响应式变量
const previewVisible = ref(false)
const previewUrlList = ref([])
const previewIndex = ref(0)

// 添加是否可以确认的计算属性
const canConfirm = computed(() => {
  if (props.multiple) {
    return selectedFiles.value.length > 0
  }
  return selectedFiles.value.length === 1
})

// 添加当前选中的分类ID
const currentCategoryId = ref(null)

// 添加移动对话框的状态
const moveDialogVisible = ref(false)
const moveTreeRef = ref(null)

// 修改上传相关的变量和方法
const uploadBtn = ref(null)

// 获取分类树数据
const getCategoryTree = async () => {
  try {
    const res = await getFileCategories()
    categoryTree.value = res || []
  } catch (error) {
    console.error('获取分类列表失败:', error)
    ElMessage.error('获取分类列表失败')
  }
}

// 获取文件列表
const getFiles = async (params = {}) => {
  try {
    const queryParams = {
      page: page.value,
      page_size: pageSize.value,  // 直接使用 pageSize
      keyword: searchKeyword.value,
      category_id: currentCategoryId.value,
      ...params // 只覆盖其他参数，不覆盖 page_size
    }
    console.log('查询参数:', queryParams)
    const res = await getFileList(queryParams)
    fileList.value = res.list || []
    total.value = res.total || 0
  } catch (error) {
    console.error('获取文件列表失败:', error)
    ElMessage.error('获取文件列表失败')
  }
}

// 处理分类点击
const handleCategoryClick = (data) => {
  console.log('点击分类:', data)
  currentCategoryId.value = data.id === 'all' ? null : Number(data.id)
  const params = {
    category_id: currentCategoryId.value,
    page: 1
  }
  getFiles(params)
}

// 处理文件选择
const handleSelectFile = (file) => {
  if (props.multiple) {
    const index = selectedFiles.value.indexOf(file.id)
    if (index > -1) {
      selectedFiles.value.splice(index, 1)
    } else {
      selectedFiles.value.push(file.id)
    }
  } else {
    selectedFiles.value = [file.id]
    handleConfirmSelect()
  }
}

// 处理搜索
const handleSearch = () => {
  page.value = 1
  getFiles() // 搜索时会带上当前分类ID
}

// 处理分页
const handleSizeChange = (val) => {
  pageSize.value = val
  getFiles() // 切换每页数量时会带上当前分类ID
}

const handleCurrentChange = (val) => {
  page.value = val
  getFiles() // 切换页码时会带上当前分类ID
}

// 处理上传
const handleUpload = () => {
  console.log('handleUpload called')
  uploadVisible.value = true
  console.log('uploadVisible:', uploadVisible.value) // 添加调试日志
}

// 自定义上传方法
const customUpload = async (options) => {
  try {
    const formData = new FormData()
    formData.append('file', options.file)
    if (currentCategoryId.value) {
      formData.append('category_id', currentCategoryId.value)
    }
    
    const res = await uploadFile(formData)
    options.onSuccess(res)
  } catch (error) {
    console.error('上传失败:', error)
    options.onError(error)
  }
}

const beforeUpload = (file) => {
  const isLt10M = file.size / 1024 / 1024 < 10
  if (!isLt10M) {
    ElMessage.error('文件大小不能超过 10MB!')
    return false
  }
  return true
}

// 修改上传成功的处理方法
const handleUploadSuccess = (response) => {
  uploadVisible.value = false
  getFiles()
}

const handleUploadError = () => {
  ElMessage.error('上传失败')
}

// 处理删除
const handleDeleteFile = (file) => {
  ElMessageBox.confirm('确定要删除该文件吗？', '提示', {
    type: 'warning'
  }).then(async () => {
    try {
      await deleteFile(file.id)
      ElMessage.success('删除成功')
      getFiles() // 刷新文件列表
    } catch (error) {
      console.error('删除失败:', error)
      ElMessage.error('删除失败')
    }
  })
}

// 批量删除
const handleRemove = () => {
  if (!selectedFiles.value.length) return
  
  ElMessageBox.confirm('确定要删除选中的文件吗？', '提示', {
    type: 'warning'
  }).then(async () => {
    try {
      await Promise.all(selectedFiles.value.map(id => deleteFile(id)))
      ElMessage.success('删除成功')
      selectedFiles.value = []
      getFiles()
    } catch (error) {
      console.error('删除失败:', error)
      ElMessage.error('删除失败')
    }
  })
}

// 工具方法
const isImage = (type) => {
  return type?.startsWith('image/')
}

const formatSize = (bytes) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const formatDate = (date) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

// 处理关闭
const handleClose = () => {
  selectedFiles.value = []
  emit('update:modelValue', false)
}

// 监听 modelValue 变化
watch(
  () => props.modelValue,
  (val) => {
    dialogVisible.value = val
    if (val && !props.standalone) {
      getCategoryTree()  // 获取分类树
      getFiles()
    }
  }
)

// 监听 dialogVisible 变化
watch(
  () => dialogVisible.value,
  (val) => {
    if (!props.standalone) {
      emit('update:modelValue', val)
    }
  }
)

// 新增功能相关的方法
const filterNode = (value, data) => {
  if (!value) return true
  return data.name.toLowerCase().includes(value.toLowerCase())
}

watch(categoryKeyword, (val) => {
  treeRef.value?.filter(val)
})

const handleSelectAll = () => {
  selectedFiles.value = fileList.value.map(file => file.id)
}

const handleInvertSelection = () => {
  const allIds = fileList.value.map(file => file.id)
  selectedFiles.value = allIds.filter(id => !selectedFiles.value.includes(id))
}

const handleClearSelection = () => {
  selectedFiles.value = []
}

const handleSort = (command) => {
  currentSort.value = command
  const sortFunctions = {
    name: (a, b) => a.name.localeCompare(b.name),
    date: (a, b) => new Date(b.createdAt) - new Date(a.createdAt),
    size: (a, b) => b.size - a.size
  }
  fileList.value.sort(sortFunctions[command])
}

const handleTableSelection = (selection) => {
  selectedFiles.value = selection.map(item => item.id)
}

const handleDownload = () => {
  const files = fileList.value.filter(file => selectedFiles.value.includes(file.id))
  files.forEach(file => {
    // 使用 fetch 下载文件
    fetch(file.url)
      .then(res => res.blob())
      .then(blob => {
        const url = window.URL.createObjectURL(blob)
        const link = document.createElement('a')
        link.href = url
        link.download = file.originalName
        document.body.appendChild(link)
        link.click()
        document.body.removeChild(link)
        window.URL.revokeObjectURL(url)
      })
      .catch(err => {
        console.error('下载失败:', err)
        ElMessage.error('下载失败')
      })
  })
}

// 单个文件下载
const handleDownloadFile = (file) => {
  // 使用 fetch 下载文件
  fetch(file.url)
    .then(res => res.blob())
    .then(blob => {
      const url = window.URL.createObjectURL(blob)
      const link = document.createElement('a')
      link.href = url
      link.download = file.originalName
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
      window.URL.revokeObjectURL(url)
    })
    .catch(err => {
      console.error('下载失败:', err)
      ElMessage.error('下载失败')
    })
}

const handleMove = () => {
  moveVisible.value = true
}

const handleMoveToCategory = async (data) => {
  if (data.id === 'all') {
    ElMessage.warning('不能移动到全部文件分类')
    return
  }

  try {
    // 批量移动文件
    await moveFiles({
      fileIds: selectedFiles.value,
      categoryId: data.id
    })
    
    ElMessage.success('移动成功')
    moveDialogVisible.value = false
    selectedFiles.value = [] // 清空选择
    // 刷新分类树数据
    getCategoryTree()
    getFiles() // 刷新文件列表
  } catch (error) {
    console.error('移动失败:', error)
    ElMessage.error('移动失败')
  }
}

// 修改图片预览处理方法
const handleImagePreview = (file) => {
  if (!isImage(file.type)) return
  
  previewUrlList.value = [file.url]
  previewIndex.value = 0
  previewVisible.value = true
}

// 修改预览方法
const handlePreview = (file) => {
  if (isImage(file.type)) {
    handleImagePreview(file)
  } else {
    // 其他类型文件的预览逻辑
    window.open(file.url)
  }
}

// 修改确认选择方法
const handleConfirmSelect = () => {
  if (!canConfirm.value) return
  
  const selectedItems = fileList.value.filter(file => 
    selectedFiles.value.includes(file.id)
  ).map(file => ({
    id: file.id,
    url: file.url,
    name: file.originalName,
    size: file.size,
    type: file.type
  }))
  
  emit('select', props.multiple ? selectedItems : selectedItems[0])
  dialogVisible.value = false
}

// 获取上传请求头
const getUploadHeaders = () => {
  try {
    const token = getToken()
    return token ? {
      Authorization: 'Bearer ' + token
    } : {}
  } catch (error) {
    console.error('获取token失败:', error)
    return {}
  }
}

// 获取上传数据
const getUploadData = () => {
  return {
    category_id: currentCategoryId.value || ''
  }
}

// 在独立页面模式下，组件挂载时初始化数据
onMounted(() => {
  if (props.standalone) {
    getCategoryTree()
    getFiles()
  }
})

// 处理粘贴事件
const handlePaste = async (e) => {
  // 阻止事件冒泡，防止重复触发
  e.stopPropagation()
  e.preventDefault()
  
  const items = e.clipboardData?.items
  if (!items) return
  
  for (let i = 0; i < items.length; i++) {
    const item = items[i]
    
    // 处理图片文件
    if (item.type.indexOf('image') !== -1) {
      const file = item.getAsFile()
      if (!file) continue
      
      // 创建上传的表单数据
      const formData = new FormData()
      formData.append('file', file)
      if (currentCategoryId.value) {
        formData.append('category_id', currentCategoryId.value)
      }
      
      try {
        // 调用上传接口
        const res = await uploadFile(formData)
        if (res) {
          ElMessage.success('粘贴上传成功')
          getFiles() // 刷新文件列表
        }
      } catch (error) {
        console.error('粘贴上传失败:', error)
        ElMessage.error('粘贴上传失败')
      }
    }
  }
}

// 处理批量移动按钮点击
const handleBatchMove = () => {
  if (!selectedFiles.value.length) return
  moveDialogVisible.value = true
}

// 处理对话框关闭
const handleDialogClose = () => {
  console.log('dialog closed')
  uploadVisible.value = false
}

// 添加监听器
watch(uploadVisible, (newVal) => {
  console.log('uploadVisible changed:', newVal)
})
</script> 