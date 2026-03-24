---
title: "Walks Are the Most Underrated Debugging Tool"
date: "2026-03-20T16:09:28+01:00"
tags: ["productivity", "programming", "health"]
description: "When I'm stuck on a hard problem, the best thing I can do is close the laptop and go outside. Science agrees."
---

Last Tuesday I spent three hours trying to fix a race condition in a concurrent pipeline. I added mutexes, rewrote the channel logic, drew diagrams. Nothing worked. Then I went for a walk around the block.

Halfway through, standing at a crosswalk, the solution appeared fully formed in my head. The bug wasn't in the concurrency at all — I was writing to a shared slice that I'd forgotten to copy. Fifteen minutes after getting back, the fix was merged.

This keeps happening.

## Why it works

Your brain has two modes of thinking. Focused mode is what you use when staring at code — linear, analytical, step-by-step. Diffuse mode is what happens when you step away — looser, more associative, better at connecting distant ideas.

The problem is that focused mode can get stuck in a loop. You keep trying the same approaches because that's all you can see. Walking forces a switch to diffuse mode, which lets your brain explore paths that focused mode had blocked.

## It's not just walking

Any low-effort physical activity works. Showering, cooking, doing dishes. The common thread is that your body is occupied with something automatic, freeing your mind to wander. The reason I prefer walking is that it also gets me outside, away from the screen, and moving.

## The hard part

The hard part isn't the walking. It's convincing yourself to leave the desk when you feel close to a solution. There's a stubborn voice that says "five more minutes" — the same voice that's been saying it for the past hour.

I've started treating it as a rule: if I've been stuck for forty-five minutes, I walk. No negotiation, no "let me try one more thing." Close the laptop, put on shoes, go.

## A broader point

We treat programming as purely intellectual work, but it happens in a body. That body needs movement, daylight, and breaks. The best code I've ever written came after stepping away from the keyboard, not after grinding through another hour at the desk.

The walk is not a break from work. The walk is part of the work.
