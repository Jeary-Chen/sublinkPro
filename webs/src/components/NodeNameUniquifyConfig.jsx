import { useState, useEffect, useMemo } from 'react';
import PropTypes from 'prop-types';

// material-ui
import Box from '@mui/material/Box';
import Collapse from '@mui/material/Collapse';
import IconButton from '@mui/material/IconButton';
import Switch from '@mui/material/Switch';
import TextField from '@mui/material/TextField';
import Typography from '@mui/material/Typography';

// icons
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import ExpandLessIcon from '@mui/icons-material/ExpandLess';
import FingerprintIcon from '@mui/icons-material/Fingerprint';

// 示例节点名用于预览
const PREVIEW_NODE_NAME = '香港节点-01';
const SECOND_PREVIEW_NODE_NAME = '日本节点';

/**
 * 节点名称唯一化配置组件
 * 通过添加机场标识前缀防止多机场间节点名称重复
 */
export default function NodeNameUniquifyConfig({ enabled, prefix, intraUniquify, airportId, onChange }) {
  const [expanded, setExpanded] = useState(false);
  const displayPrefix = prefix.trim() || (airportId ? `[A${airportId}]` : '[A保存后分配ID]');

  const getSettingRowSx = (active) => ({
    px: 1.5,
    py: 1,
    borderRadius: 1.5,
    border: '1px solid',
    borderColor: 'divider',
    bgcolor: active ? 'background.default' : 'action.disabledBackground',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'space-between',
    gap: 1.5
  });

  // 根据当前配置计算预览结果
  const previewResult = useMemo(() => {
    if (!enabled) {
      return PREVIEW_NODE_NAME;
    }
    return displayPrefix + PREVIEW_NODE_NAME;
  }, [enabled, displayPrefix]);

  const duplicatePreviewGroups = useMemo(() => {
    if (!intraUniquify) {
      return [];
    }

    const previewPrefix = enabled ? displayPrefix : '';

    return [
      {
        label: `${PREVIEW_NODE_NAME} 重名组`,
        items: [`${previewPrefix}${PREVIEW_NODE_NAME}-1`, `${previewPrefix}${PREVIEW_NODE_NAME}-2`]
      },
      {
        label: `${SECOND_PREVIEW_NODE_NAME} 重名组`,
        items: [`${previewPrefix}${SECOND_PREVIEW_NODE_NAME}-1`, `${previewPrefix}${SECOND_PREVIEW_NODE_NAME}-2`]
      }
    ];
  }, [enabled, intraUniquify, displayPrefix]);

  // 自动展开面板（当开启功能时）
  useEffect(() => {
    if (enabled || intraUniquify) {
      setExpanded(true);
    }
  }, [enabled, intraUniquify]);

  const handleEnabledChange = (event) => {
    onChange({ enabled: event.target.checked, prefix, intraUniquify });
  };

  const handlePrefixChange = (event) => {
    onChange({ enabled, prefix: event.target.value, intraUniquify });
  };

  const handleIntraUniquifyChange = (event) => {
    onChange({ enabled, prefix, intraUniquify: event.target.checked });
  };

  return (
    <Box
      sx={{
        border: '1px solid',
        borderColor: 'divider',
        borderRadius: 2,
        overflow: 'hidden'
      }}
    >
      {/* 标题栏 */}
      <Box
        sx={{
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between',
          px: 2,
          py: 1,
          bgcolor: 'background.default',
          cursor: 'pointer',
          '&:hover': { bgcolor: 'action.hover' }
        }}
        onClick={() => setExpanded(!expanded)}
      >
        <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'flex-start', gap: 0.75, minWidth: 0, flex: 1 }}>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, minWidth: 0, width: '100%' }}>
            <FingerprintIcon color="action" fontSize="small" />
            <Box sx={{ minWidth: 0 }}>
              <Typography variant="body2" fontWeight={500}>
                节点名称唯一化
              </Typography>
              <Typography
                variant="caption"
                sx={{
                  display: 'block',
                  color: 'text.secondary'
                }}
              >
                跨机场前缀与机场内重名编号可分别开启，互不影响
              </Typography>
            </Box>
          </Box>
          {(enabled || intraUniquify) && (
            <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.75 }}>
              {enabled && (
                <Typography
                  variant="caption"
                  sx={{
                    px: 1,
                    py: 0.25,
                    borderRadius: 1,
                    bgcolor: 'primary.main',
                    color: 'primary.contrastText'
                  }}
                >
                  跨机场前缀
                </Typography>
              )}
              {intraUniquify && (
                <Typography
                  variant="caption"
                  sx={{
                    px: 1,
                    py: 0.25,
                    borderRadius: 1,
                    bgcolor: 'secondary.main',
                    color: 'secondary.contrastText'
                  }}
                >
                  机场内编号
                </Typography>
              )}
            </Box>
          )}
        </Box>
        <IconButton size="small">{expanded ? <ExpandLessIcon fontSize="small" /> : <ExpandMoreIcon fontSize="small" />}</IconButton>
      </Box>

      {/* 展开内容 */}
      <Collapse in={expanded}>
        <Box sx={{ p: 2, borderTop: '1px solid', borderColor: 'divider' }}>
          <Typography variant="caption" color="textSecondary" sx={{ display: 'block', mb: 2 }}>
            跨机场前缀唯一化与机场内重名编号是两个独立设置：前者用于区分不同机场的同名节点，后者仅处理同一机场内的重名节点。
          </Typography>

          <Box sx={{ display: 'grid', gap: 1.25, mb: 2 }}>
            <Box sx={{ ...getSettingRowSx(enabled), display: 'grid', gap: 1.25 }}>
              <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', gap: 1.5 }}>
                <Box sx={{ minWidth: 0 }}>
                  <Typography variant="body2" fontWeight={500}>
                    跨机场前缀
                  </Typography>
                  <Typography variant="caption" color="text.secondary">
                    为当前机场节点添加唯一前缀，用于区分不同机场之间的同名节点
                  </Typography>
                </Box>
                <Switch size="small" checked={enabled} onChange={handleEnabledChange} />
              </Box>
              {enabled && (
                <TextField
                  fullWidth
                  size="small"
                  label="跨机场前缀（可选）"
                  placeholder={airportId ? `留空使用默认前缀 [A${airportId}]` : '留空后保存时使用默认前缀 [A新机场ID]'}
                  value={prefix}
                  onChange={handlePrefixChange}
                  helperText="仅在开启跨机场前缀唯一化时生效，用于区分不同机场的同名节点"
                />
              )}
            </Box>

            <Box sx={getSettingRowSx(intraUniquify)}>
              <Box sx={{ minWidth: 0 }}>
                <Typography variant="body2" fontWeight={500}>
                  机场内重名编号
                </Typography>
                <Typography variant="caption" color="text.secondary">
                  同一机场内若有重名节点，会按各自重名组在名称后追加 -1、-2、-3…；不同名称组会分别从 -1 开始
                </Typography>
              </Box>
              <Switch size="small" checked={intraUniquify} onChange={handleIntraUniquifyChange} />
            </Box>
          </Box>

          {/* 预览效果 */}
          <Box
            sx={{
              p: 1.5,
              borderRadius: 1,
              bgcolor: enabled || intraUniquify ? 'action.selected' : 'action.disabledBackground'
            }}
          >
            <Typography variant="caption" color="textSecondary" sx={{ display: 'block', mb: 0.5 }}>
              跨机场前缀预览
            </Typography>
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, flexWrap: 'wrap' }}>
              <Typography variant="body2" sx={{ fontFamily: 'monospace', color: 'text.secondary', textDecoration: 'line-through' }}>
                {PREVIEW_NODE_NAME}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                →
              </Typography>
              <Typography
                variant="body2"
                sx={{
                  fontFamily: 'monospace',
                  fontWeight: 500,
                  color: enabled ? 'primary.main' : 'text.disabled'
                }}
              >
                {previewResult}
              </Typography>
            </Box>

            <Typography variant="caption" color="textSecondary" sx={{ display: 'block', mt: 0.75 }}>
              {enabled ? '开启后会为当前机场添加唯一前缀，仅影响不同机场之间的同名节点。' : '未开启时，跨机场同名节点仍保持原名称。'}
            </Typography>

            {intraUniquify && (
              <Box sx={{ mt: 1.5 }}>
                <Typography variant="caption" color="textSecondary" sx={{ display: 'block', mb: 0.5 }}>
                  机场内重名编号预览
                </Typography>
                <Typography variant="caption" color="textSecondary" sx={{ display: 'block', mb: 1 }}>
                  每组重名会单独编号，HK 与日本两个名称组都会各自从 -1 开始，而不是共用一套全局序号。
                </Typography>

                {duplicatePreviewGroups.map((group) => (
                  <Box key={group.label} sx={{ '&:not(:last-of-type)': { mb: 1.25 } }}>
                    <Typography variant="caption" color="text.secondary" sx={{ display: 'block', mb: 0.5 }}>
                      {group.label}
                    </Typography>
                    {group.items.map((item) => (
                      <Typography key={item} variant="body2" sx={{ fontFamily: 'monospace', color: 'primary.main', fontWeight: 500 }}>
                        {item}
                      </Typography>
                    ))}
                  </Box>
                ))}
              </Box>
            )}
          </Box>
        </Box>
      </Collapse>
    </Box>
  );
}

NodeNameUniquifyConfig.propTypes = {
  enabled: PropTypes.bool.isRequired,
  prefix: PropTypes.string,
  intraUniquify: PropTypes.bool,
  airportId: PropTypes.number,
  onChange: PropTypes.func.isRequired
};

NodeNameUniquifyConfig.defaultProps = {
  prefix: '',
  intraUniquify: false,
  airportId: 0
};
