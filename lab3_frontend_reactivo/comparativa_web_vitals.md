# Auditoría de Rendimiento y Web Vitals: Next.js 15 / React 19 (VDOM) vs. Svelte 5 Runes (AOT)

## 1. Presupuesto de Cuadro (*Frame Budget*) — $60\text{ Hz}$ y $120\text{ Hz}$

El tiempo disponible por cuadro (*frame*) es el inverso de la frecuencia de refresco de pantalla ($f$):

$$T_{\text{frame}} = \frac{1000\text{ ms}}{f\text{ (Hz)}}$$

* **A $60\text{ Hz}$:**
  $$T_{60\text{ Hz}} = \frac{1000\text{ ms}}{60} \approx 16.67\text{ ms}$$
* **A $120\text{ Hz}$:**
  $$T_{120\text{ Hz}} = \frac{1000\text{ ms}}{120} \approx 8.33\text{ ms}$$

Considerando el margen reservado por el navegador para *layout*, *paint*, *compositing* y recolección de basura/eventos del sistema (~$4\text{–}6\text{ ms}$), el cómputo de *scripts* y actualización del DOM en una ráfaga de telemetría 5G debería mantenerse idealmente por debajo de **$\sim 10\text{ ms}$ a $60\text{ Hz}$** y por debajo de **$\sim 4\text{ ms}$ a $120\text{ Hz}$** para evitar *frame drops* (*jank*).

## 2. Tabla Comparativa de Arquitecturas

| Métrica / Criterio | Next.js 15 + React 19 (`TelemetryPanelClient.tsx`) | Svelte 5 Runes AOT (`+page.svelte`) |
| :--- | :--- | :--- |
| **Modelo de reactividad y complejidad** | Reconciliación con Virtual DOM (`useState` / `useMemo`): recorrido de subárboles y *diffing* proporcional al árbol afectado, $O(N)$ en nodos evaluados. | Grafo acíclico dirigido (DAG) de señales finas compilado AOT (`$state` / `$derived.by`): actualización quirúrgica directa sobre nodos DOM dependientes, $O(1)$ por señal afectada. |
| **Tamaño del bundle JS (cliente)** | **Mayor**: incluye runtime de React, reconciliador VDOM y *hooks* además del código del componente. | **Menor**: el compilador AOT emite instrucciones directas de actualización de DOM sin runtime de VDOM. |
| **Tiempo de hidratación** | **Moderado / Mayor**: requiere rehidratar el árbol de componentes cliente tras la frontera `"use client"` y enlazar *hooks* y *handlers*. | **Bajo / Rápido**: *hydration* ligera y granular gracias a la compilación anticipada (AOT) y reactividad fina. |
| **INP (*Interaction to Next Paint*)** | Sensible a ráfagas grandes y re-renders en cascada si el trabajo de reconciliación supera el *frame budget* ($16.67\text{ ms}$ / $8.33\text{ ms}$). | Muy competitivo en ráfagas frecuentes al evitar *diffing* de VDOM y tocar únicamente los bindings vinculados a cada señal. |