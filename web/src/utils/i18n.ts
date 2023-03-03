// translate router.meta.title, be used in breadcrumb sidebar
import i18n from '@/lang/index';

export function generateTitle(title: any) {
  // 判断是否存在国际化配置，如果没有原生返回
  console.log('route.' + title)
  const hasKey = i18n.global.te('route.' + title);
  if (hasKey) {
    const translatedTitle = i18n.global.t('route.' + title);
    return translatedTitle;
  }
  return title;
}

export function t(v: any) {
  const hasKey = i18n.global.te(v);
  if (hasKey) {
    return i18n.global.t(v);
  }
  return v;
}
