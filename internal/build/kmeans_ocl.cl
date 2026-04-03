kernel void assign_first_iter(
    read_only image2d_t inputImg,
    sampler_t smp,
    read_only image2d_t clusters,
    __global ushort* clusterAssignments,
    write_only image2d_t sampleResult)
{
    uint x = get_global_id(0);
    uint y = get_global_id(1);

    uint inpW = get_image_width(inputImg);
    uint inpH = get_image_height(inputImg);

    uint dstW = get_image_width(sampleResult);
    uint dstH = get_image_height(sampleResult);

    if (x >= dstW || y >= dstH) {
        return;
    }

    uint k = get_image_width(clusters);

    float4 pixel;
    if (inpW == dstW && inpH == dstH) {
        pixel = read_imagef(inputImg, (int2)(x, y));
    } else {
        float2 normCoord = (float2)(
            (float)x / (float)dstW,
            (float)y / (float)dstH
        );
        pixel = read_imagef(inputImg, smp, normCoord);
    }
    write_imagef(sampleResult, (int2)(x, y), pixel);

    float minDistance = FLT_MAX;
    ushort assign = USHRT_MAX;
    for (int i = 0; i < k; i++) {
        float4 cluster = read_imagef(clusters, (int2)(i, 0));
        float4 diff = pixel - cluster;
        diff.w = 0;
        float d = dot(diff, diff);
        if (d < minDistance) {
            minDistance = d;
            assign = (ushort)i;
        }
    }
    clusterAssignments[x+y*dstW] = assign;
}

kernel void assign_remaining_iter(
    read_only image2d_t sampleResult,
    read_only image2d_t clusters,
    __global ushort* clusterAssignments)
{
    uint x = get_global_id(0);
    uint y = get_global_id(1);

    uint dstW = get_image_width(sampleResult);
    uint dstH = get_image_height(sampleResult);

    if (x >= dstW || y >= dstH) {
        return;
    }

    uint k = get_image_width(clusters);

    float4 pixel = read_imagef(sampleResult, (int2)(x, y));

    float minDistance = FLT_MAX;
    ushort assign = USHRT_MAX;
    for (int i = 0; i < k; i++) {
        float4 cluster = read_imagef(clusters, (int2)(i, 0));
        float4 diff = pixel - cluster;
        diff.w = 0;
        float d = dot(diff, diff);
        if (d < minDistance) {
            minDistance = d;
            assign = (ushort)i;
        }
    }
    clusterAssignments[x+y*dstW] = assign;
}
